//usr/local/go/bin/go run $0 $@ $(dirname `realpath $0`); exit
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

type depSpec struct {
	url    string
	branch string
	commit string
	path   string
	isGit  bool
}

type buildConf struct {
	PRE_BUILD_EXEC []string
	PRE_DEBUG_EXEC []string
	POST_COMP_EXEC []string
	EXTRA_LDFLAGS  string
	EXTRA_ENVIRON  []string
	BUILD_TAGS     string
}

var bc buildConf

var (
	PROJ_ROOT, PROJ_NAME, CMD string
	PROJ_ARGS                 []string
)

func getDepends() (deps []depSpec) {
	f, err := os.Open("depends")
	if os.IsNotExist(err) {
		return
	}
	assert(err)
	defer f.Close()
	spec := bufio.NewScanner(f)
	for spec.Scan() {
		s := strings.TrimSpace(spec.Text())
		if len(s) == 0 || strings.HasPrefix(s, "#") {
			continue
		}
		spec := strings.Split(s, " ")
		if len(spec) > 2 {
			panic(fmt.Errorf("invalid spec: %s", s))
		}
		var gs depSpec
		gs.isGit = true
		gs.url = spec[0]
		u, err := url.Parse(gs.url)
		if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
			gs.path = path.Join(u.Host, u.Path)
		} else {
			gs.path = gs.url
			u := strings.SplitN(gs.url, "@", 2)
			if len(u) == 2 {
				gs.path = strings.Replace(u[1], ":", "/", -1)
			} else {
				gs.isGit = false
			}
		}
		if strings.HasSuffix(gs.path, ".git") {
			gs.path = gs.path[:len(gs.path)-4]
		}
		gs.commit = "HEAD"
		gs.branch = "master"
		if len(spec) == 2 {
			if !gs.isGit {
				panic(fmt.Errorf("cannot specify branch/commit for 'go get' spec"))
			}
			bnc := strings.SplitN(spec[1], "@", 2)
			if bnc[0] != "" {
				gs.branch = bnc[0]
			}
			if len(bnc) == 2 {
				gs.commit = bnc[1]
			}
		}
		deps = append(deps, gs)
	}
	return
}

func updDepends(deps []depSpec, mode int) (depRoots []string) {
	PROJ_SRC := path.Join(PROJ_ROOT, "src")
	rs := make(map[string]int)
	for _, repo := range deps {
		fmt.Printf("clone: %s %s@%s", repo.url, repo.branch, repo.commit)
		root := strings.SplitN(repo.path, "/", 2)[0]
		rs[root] = 1
		cd := path.Join(PROJ_SRC, repo.path)
		fi, err := os.Stat(path.Join(cd, ".git"))
		if err == nil && fi.IsDir() {
			switch mode {
			case 1: //force re-sync
				err := exec.Command("rm", "-fr", cd).Run()
				if err != nil {
					fmt.Printf("\nFAILED: rm -fr %s (%v)\n", cd, err)
					os.Exit(1)
				}
			case 2: //update repo
				if repo.isGit {
					_, err = exec.Command("git", "-C", cd, "pull").Output()
					assert(err)
				} else {
					assert(run("go", "get", "-u", repo.url))
				}
				fmt.Println(" ...updated")
				continue
			default: //skip if exists
				fmt.Println(" ...skipped")
				continue
			}
		}
		fmt.Printf("\n")
		if !repo.isGit {
			assert(run("go", "get", repo.url))
			continue
		}
		args := []string{"clone", repo.url}
		if repo.commit == "HEAD" {
			args = append(args, "--depth", "1")
		}
		args = append(args, "--branch", repo.branch, "--single-branch", cd)
		_, err = exec.Command("git", args...).Output()
		if err != nil {
			switch err.(type) {
			case *exec.ExitError:
				fmt.Println(string(err.(*exec.ExitError).Stderr))
			default:
				fmt.Println(err)
			}
			os.Exit(1)
		}
		if repo.commit != "HEAD" {
			os.Chdir(cd)
			assert(exec.Command("git", "checkout", repo.commit).Run())
		}
	}
	for r := range rs {
		depRoots = append(depRoots, r)
	}
	return
}

func updGitIgnore(roots []string) {
	patterns := make(map[string]int)
	buf, err := ioutil.ReadFile(".gitignore")
	if os.IsNotExist(err) {
		f, err := os.Create(".gitignore")
		assert(err)
		defer f.Close()
		for _, p := range roots {
			f.WriteString(p + "\n")
		}
		return
	}
	assert(err)
	for _, p := range strings.Split(string(buf), "\n") {
		patterns[p] = 1
	}
	for _, p := range roots {
		patterns[p] = 1
	}
	f, err := os.Create(".gitignore")
	assert(err)
	defer f.Close()
	var ps []string
	for p := range patterns {
		p = strings.TrimSpace(p)
		if p != "" {
			ps = append(ps, p)
		}
	}
	sort.Strings(ps)
	for _, p := range ps {
		f.WriteString(p + "\n")
	}
}

func getGitInfo() (branch, hash string, revisions int) {
	branch = "unknown"
	hash = "unknown"
	o, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return
	}
	branch = strings.TrimSpace(string(o))
	o, err = exec.Command("git", "log", "-n1", "--pretty=format:%h").Output()
	if err != nil {
		return
	}
	hash = string(o)
	o, err = exec.Command("git", "log", "--oneline").Output()
	if err != nil {
		return
	}
	revisions = len(strings.Split(string(o), "\n")) - 1
	return
}

func parseConf() (err error) {
	conf := path.Join(PROJ_ROOT, "src", PROJ_NAME, "build.conf")
	f, err := os.Open(conf)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("[build.conf] not found, continue with defaults...")
			err = nil
		}
		return
	}
	defer f.Close()
	getCmd := func(cmdline string) []string {
		cmdline = strings.TrimSpace(cmdline)
		w := "/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789."
		if strings.ContainsAny(cmdline[0:1], w) {
			return strings.Split(cmdline, " ")
		}
		return strings.Split(cmdline[1:], cmdline[0:1])
	}
	lines := bufio.NewScanner(f)
	for lines.Scan() {
		line := strings.TrimSpace(lines.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			continue
		}
		key := strings.TrimSpace(kv[0])
		switch strings.ToUpper(key) {
		case "PRE_BUILD_EXEC":
			bc.PRE_BUILD_EXEC = getCmd(kv[1])
		case "PRE_DEBUG_EXEC":
			bc.PRE_DEBUG_EXEC = getCmd(kv[1])
		case "POST_COMP_EXEC":
			bc.POST_COMP_EXEC = getCmd(kv[1])
		case "EXTRA_LDFLAGS":
			bc.EXTRA_LDFLAGS = strings.TrimSpace(kv[1])
		case "EXTRA_ENVIRON":
			bc.EXTRA_ENVIRON = getCmd(kv[1])
		case "BUILD_TAGS":
			bc.BUILD_TAGS = strings.TrimSpace(kv[1])
		default:
			err = fmt.Errorf("Invalid configuration key: %s", key)
			return
		}
	}
	err = lines.Err()
	return
}

func run(args ...string) (err error) {
	cmd := exec.Command(args[0], args[1:]...)
	for _, e := range os.Environ() {
		cmd.Env = append(cmd.Env, e)
	}
	cmd.Env = append(cmd.Env, "GOPATH="+PROJ_ROOT)
	if args[0] == "go" && args[1] == "build" {
		if CMD == "win" {
			cmd.Env = append(cmd.Env, "GOOS=windows")
		} else if CMD == "arm" {
			cmd.Env = append(cmd.Env, "GOARCH=arm", "GOARM=7")
		}
		cmd.Env = append(cmd.Env, bc.EXTRA_ENVIRON...)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func parseCmdline() {
	PROJ_ROOT = os.Args[len(os.Args)-1]
	if CMD == "sync" {
		return
	}
	var mains, main []string
	filepath.Walk(PROJ_ROOT, func(path string, info os.FileInfo,
		err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			func() {
				dir := filepath.Dir(path[len(PROJ_ROOT)+5:])
				if dir == "." {
					return
				}
				f, err := os.Open(path)
				if err != nil {
					return
				}
				defer f.Close()
				score := 0
				lines := bufio.NewScanner(f)
				for lines.Scan() {
					line := lines.Text()
					if strings.HasPrefix(line, "package main") {
						score++
					}
					if strings.HasPrefix(line, "func main()") {
						score++
					}
					if score > 1 {
						break
					}
				}
				if score > 1 {
					mains = append(mains, dir)
				}
			}()
		}
		return nil
	})
	if len(mains) == 0 {
		fmt.Println("No target found (require a func main() in package main)")
		os.Exit(1)
	}
	PROJ_NAME = path.Base(PROJ_ROOT)
	PROJ_ARGS = os.Args[1 : len(os.Args)-1]
	if len(PROJ_ARGS) > 0 && !strings.HasPrefix(PROJ_ARGS[0], "-") {
		PROJ_NAME = PROJ_ARGS[0]
		PROJ_ARGS = PROJ_ARGS[1:]
	}
	for _, m := range mains {
		if m == PROJ_NAME {
			main = []string{m}
			break
		}
		if strings.Contains(m, PROJ_NAME) {
			main = append(main, m)
		}
	}
	switch len(main) {
	case 0:
		fmt.Printf("Invalid sub-project [%s], valid names:", PROJ_NAME)
		for _, m := range mains {
			fmt.Printf(" [%s]", m)
		}
		fmt.Println()
		os.Exit(1)
	case 1:
		PROJ_NAME = main[0]
	default:
		fmt.Print("Ambiguous sub-project name, matched:")
		for _, m := range main {
			fmt.Printf(" [%s]", m)
		}
		fmt.Println()
		os.Exit(1)
	}
}

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("PANIC:", e.(error).Error())
			os.Exit(1)
		}
	}()
	CMD = path.Base(os.Args[0])
	parseCmdline()
	assert(os.Chdir(PROJ_ROOT))
	var syncMode int
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-force":
			syncMode = 1
		case "-update":
			syncMode = 2
		}
	}
	depRoots := updDepends(getDepends(), syncMode)
	updGitIgnore(depRoots)
	if CMD == "sync" {
		return
	}
	fmt.Println()
	err := parseConf()
	if err != nil {
		fmt.Printf("parseConf: %s\n", err)
		os.Exit(1)
	}
	scripts := bc.PRE_BUILD_EXEC
	if CMD == "run" && len(bc.PRE_DEBUG_EXEC) > 0 {
		scripts = bc.PRE_DEBUG_EXEC
	}
	if len(scripts) > 0 {
		fmt.Printf("PRE_COMP_EXEC: %s\n", strings.Join(scripts, " "))
		err = run(scripts...)
		if err != nil {
			fmt.Printf("PRE_COMP_EXEC: %s\n", err)
			os.Exit(1)
		}
	}
	branch, hash, revs := getGitInfo()
	args := []string{"go"}
	if CMD != "build" && CMD != "run" {
		exe := PROJ_NAME
		if CMD == "win" {
			exe += ".exe"
		}
		args = append(args, "build", "-o", "bin/"+exe)
	} else {
		args = append(args, "install")
	}
	//add race condiction detection for debugging purpose
	if CMD == "run" {
		args = append(args, "-race")
	}
	if len(bc.BUILD_TAGS) > 0 {
		args = append(args, "-tags", bc.BUILD_TAGS)
	}
	ldflags := fmt.Sprintf(`%s -s -w -X main._G_BRANCH=%s -X main._G_HASH=%s
		-X main._G_REVS=%d -X main._BUILT_=%d`, bc.EXTRA_LDFLAGS, branch,
		hash, revs, time.Now().Unix())
	args = append(args, "-ldflags", ldflags, PROJ_NAME)
	err = run(args...)
	if err != nil {
		fmt.Printf("COMPILE: %s\n", err)
		os.Exit(1)
	}
	if len(bc.POST_COMP_EXEC) > 0 {
		fmt.Printf("POST_COMP_EXEC: %s\n", strings.Join(bc.POST_COMP_EXEC, " "))
		err = run(bc.POST_COMP_EXEC...)
		if err != nil {
			fmt.Printf("POST_COMP_EXEC: %s", err)
			os.Exit(1)
		}
	}
	if CMD == "run" {
		fmt.Println("\nRUNNING:")
		args := []string{path.Join(PROJ_ROOT, "bin", PROJ_NAME)}
		args = append(args, PROJ_ARGS...)
		assert(run(args...))
	}
}
