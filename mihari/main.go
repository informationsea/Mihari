package main

import ("fmt"
	"os"
	"runtime"
	"os/exec"
	"io"
	"io/ioutil"
	"path"
	"bufio"
	"log"
	"time"
)

func main() {

	// Check argument

	if len(os.Args) == 1 {
		fmt.Printf("Mihari\n"+
			"mihari commands...\n"+
			"mihari --init\n" +
			"mihari --makefile\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--init":
		initMode(true)
	case "--makefile":
		makefileMode()
	default:
		loggingMode()
	}
	
}

func loggingMode() {
	// search log dir
	logdir, err := searchLogdir()
	if err != nil {
		initMode(false)
	}
	fmt.Printf("logdir: %s\n", logdir)

	// extract assets
	lib, err := os.Create(path.Join(logdir, "libmiharihook.so"))
	if err != nil {log.Fatal(err); os.Exit(1)}
	libdata, err := Asset("assets/libmiharihook.so")
	if err != nil {log.Fatal(err); os.Exit(1)}
	_, err = lib.Write(libdata)
	if err != nil {log.Fatal(err); os.Exit(1)}
	lib.Close()

	// setup environment
	if runtime.GOOS == "darwin" {
		os.Setenv("DYLD_INSERT_LIBRARIES", path.Join(logdir, "libmiharihook.so"))
		os.Setenv("DYLD_FORCE_FLAT_NAMESPACE", "YES")
	} else if runtime.GOOS == "linux" {
		os.Setenv("LD_PRELAOD", "...")
	}

	// setup and start process
	logfile, err := ioutil.TempFile("", "mihrai-log-")
	if err != nil {log.Fatal(err); os.Exit(1)}
	os.Setenv("MIHARI_LOG_FILE", logfile.Name())
	//fmt.Printf("logfile: %s\n", logfile.Name())
	fmt.Printf("Command: %s\n", escapeCommand(os.Args[1:]))

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {log.Fatal(err); os.Exit(1)}
	stderr, err := cmd.StderrPipe()
	if err != nil {log.Fatal(err); os.Exit(1)}
	err = cmd.Start()
	if err != nil {log.Fatal(err); os.Exit(1)}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	
	cmd.Wait()

	baseDir := path.Dir(logdir)
	runlog, err := NewRunLogFromExecuteLog(os.Args[1:], time.Now(), logfile, baseDir)
	if err != nil {log.Fatal(err); os.Exit(1)}
	runlogFile, err := os.Create(path.Join(logdir, "log-"+NormalizeArgs(os.Args[1:])))
	if err != nil {log.Fatal(err); os.Exit(1)}
	WriteRunLog(runlogFile, runlog)

	os.Remove(logfile.Name())
}

func initMode(createWithouQuestion bool) {
	logdir, err := searchLogdir()
	if err != nil {
		if createWithouQuestion {
			currentDir, err := os.Getwd()
			logdir = path.Join(currentDir, ".mihari")
			err = os.Mkdir(logdir, 0755)
			if err != nil {log.Fatal(err); os.Exit(1)}
		} else {
			fmt.Println("Mihari logging directory is not found.")
			fmt.Println("Do you want to create it? (yes/no)")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				switch scanner.Text() {
				case "yes":
					currentDir, err := os.Getwd()
					logdir = path.Join(currentDir, ".mihari")
					err = os.Mkdir(logdir, 0755)
					if err != nil {log.Fatal(err); os.Exit(1)}
					goto exitlogdirquestion
				case "no":
					fmt.Println("No logging dir")
					os.Exit(1)
				default:
					fmt.Println("yes or no")
				}
			}
		}
	} else {
		fmt.Println("Mihari logging directory is present")
	}

exitlogdirquestion:

}

func makefileMode() {
	fmt.Println("makefile mode")
	
	logdir, err := searchLogdir()
	if err != nil {
		log.Fatal(err)
		fmt.Println("No loggin directory")
		os.Exit(1)
	}

	logdirFile, err := os.Open(logdir)
	if err != nil {log.Fatal(err); os.Exit(1)}
	fi, err := logdirFile.Readdir(0)

	runLogList := []RunLog{}

	for _, one := range(fi) {
		if !one.IsDir() && one.Name()[:4] == "log-" {
			oneLogPath, err := os.Open(path.Join(logdir, one.Name()))
			if err != nil {log.Fatal(err); os.Exit(1)}
			oneLog, err := LoadRunLog(oneLogPath)
			if err != nil {log.Fatal(err); os.Exit(1)}	
			runLogList= append(runLogList, oneLog)
		}
	}

	makefile, err := os.Create("./Makefile")
	if err != nil {log.Fatal(err); os.Exit(1)}	
	files, commands, err := ResolveDependency(runLogList)
	if err != nil {log.Fatal(err); os.Exit(1)}	
	err = generateMakefile(files, commands, makefile)
	if err != nil {log.Fatal(err); os.Exit(1)}
	makefile.Close()
}
