package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/kardianos/service"
	"github.com/topxeq/tk"
)

var versionG = "0.9a"
var serviceNameG = "tasksX"
var basePathG = ""
var configFileNameG = serviceNameG + ".cfg"
var serviceModeG = false
var runModeG = ""
var currentOSG = ""
var currentPortG = "7489"

type program struct {
	BasePath string
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	// basePathG = p.BasePath
	// logWithTime("basePath: %v", basePathG)
	serviceModeG = true

	go p.run()

	return nil
}

func (p *program) run() {
	go doWork()
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func plByMode(formatA string, argsA ...interface{}) {
	if runModeG == "cmd" {
		tk.Pl(formatA, argsA...)
	} else {
		tk.AddDebugF(formatA, argsA...)
	}
}

func initSvc() *service.Service {
	svcConfigT := &service.Config{
		Name:        serviceNameG,
		DisplayName: serviceNameG,
		Description: serviceNameG + " V" + versionG,
	}

	prgT := &program{BasePath: basePathG}
	var s, err = service.New(prgT, svcConfigT)

	if err != nil {
		tk.LogWithTimeCompact("%s unable to start: %s\n", svcConfigT.DisplayName, err)
		return nil
	}

	return &s
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	if req != nil {
		req.ParseForm()
	}

	// reqT := tk.GetFormValueWithDefaultValue(req, "prms", "")

	plByMode("req: %+v", req)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Test."))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func doJapi(res http.ResponseWriter, req *http.Request) string {

	defer func() {
		if r := recover(); r != nil {
			tk.AddDebugF("japi: Recovered: %v\n%v", r, string(debug.Stack()))
			tk.AddDebugF("japi Recovered: %v\n%v", r, string(debug.Stack()))
		}
	}()

	if req != nil {
		req.ParseForm()
	}

	reqT := tk.GetFormValueWithDefaultValue(req, "req", "")

	if res != nil {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "*")
		res.Header().Set("Content-Type", "text/json;charset=utf-8")
	}

	res.WriteHeader(http.StatusOK)

	switch reqT {

	case "debug":
		{
			tk.Pl("%v", req)
			a := make([]int, 3)
			a[5] = 8

			return tk.GenerateJSONPResponse("success", tk.IntToStr(a[5]), req)
		}

	case "getDebug":
		{
			res.Header().Set("Content-Type", "text/plain;charset=utf-8")

			res.WriteHeader(http.StatusOK)

			return tk.GenerateJSONPResponse("success", tk.GetDebug(), req)
		}

	case "clearDebug":
		{
			tk.ClearDebug()
			return tk.GenerateJSONPResponse("success", "", req)
		}

	case "requestinfo":
		{
			rs := tk.Spr("%#v", req)

			return tk.GenerateJSONPResponse("success", rs, req)
		}
	default:
		return tk.GenerateJSONPResponse("fail", tk.Spr("unknown request: %v", req), req)
	}
}

func japiHandler(w http.ResponseWriter, req *http.Request) {
	rs := doJapi(w, req)

	// w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte(rs))
}

var staticFS http.Handler = nil
var staticTemplateFS http.Handler = nil

func serveStaticDirHandler(w http.ResponseWriter, r *http.Request) {
	plByMode("in serveStaticDir")
	if staticFS == nil {
		staticFS = http.StripPrefix("/w/", http.FileServer(http.Dir(filepath.Join(basePathG, "w"))))
	}

	old := r.URL.Path

	// tk.Pl("urlPath: %v", r.URL.Path)

	name := filepath.Join(basePathG, path.Clean(old))

	// tk.Pl("name: %v", name)

	info, err := os.Lstat(name)
	if err == nil {
		if !info.IsDir() {
			staticFS.ServeHTTP(w, r)
			// http.ServeFile(w, r, name)
		} else {
			if tk.IfFileExists(filepath.Join(name, "index.html")) {
				staticFS.ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		}
	} else {
		http.NotFound(w, r)
	}

}

func doHttp(res http.ResponseWriter, req *http.Request) string {

	defer func() {
		if r := recover(); r != nil {
			tk.AddDebugF("http Recovered: %v\n%v", r, string(debug.Stack()))
			tk.AddDebugF("http Recovered: %v\n%v", r, string(debug.Stack()))
		}
	}()

	if req != nil {
		req.ParseForm()
	}

	reqT := tk.GetFormValueWithDefaultValue(req, "req", "")

	if res != nil {
		res.Header().Set("Access-Control-Allow-Origin", "*")
	}

	if reqT == "" {
		if tk.StartsWith(req.RequestURI, "/dp") {
			reqT = req.RequestURI[3:]
		}
	}

	tmps := tk.Split(reqT, "?")
	if len(tmps) > 1 {
		reqT = tmps[0]
	}

	plByMode("reqT: %v, req: %+v", reqT, req)

	switch reqT {

	case "test":
		{
			return tk.GenerateJSONPResponse("success", "test", req)
		}
	case "qr", "qr/", "/qr", "/qr/":
		{
			// contentT := tk.GetFormValueWithDefaultValue(req, "content", "")
			// if contentT == "" {
			// 	contentT = "http://topget.org"
			// }

			// qrCode, _ := qr.Encode(contentT, qr.M, qr.Auto)

			// // Scale the barcode to 200x200 pixels
			// qrCode, _ = barcode.Scale(qrCode, 500, 500)

			// // // create the output file
			// // file, _ := os.Create("qrcode.png")
			// // defer file.Close()

			// // // encode the barcode as png
			// // png.Encode(file, qrCode)

			// res.Header().Set("Content-Type", "image/png")

			// png.Encode(res, qrCode)

			return "TX_END_RESPONSE_XT"
		}
	default:
		{
			return ""
		}
	}
}

func httpHandler(w http.ResponseWriter, req *http.Request) {
	plByMode("in httpHandler")

	rs := doHttp(w, req)

	if rs == "TX_END_RESPONSE_XT" {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(rs))
}

func startHttpServer(portA string) {
	defer func() {
		if r := recover(); r != nil {
			tk.LogWithTimeCompact("startHttpServer: Recovered: %v\n%v", r, string(debug.Stack()))
		}
	}()

	tk.LogWithTimeCompact("trying startHttpServer, port: %v", portA)

	http.HandleFunc("/w/", serveStaticDirHandler)

	http.HandleFunc("/dp/", httpHandler)

	http.HandleFunc("/japi", japiHandler)

	http.HandleFunc("/", mainHandler)

	err := http.ListenAndServe(":"+portA, nil)
	if err != nil {
		plByMode("ListenAndServeHttp: %v", err.Error())
		tk.LogWithTimeCompact("ListenAndServeHttp: %v", err.Error())
	} else {
		plByMode("ListenAndServeHttp: %v", currentPortG)
		tk.LogWithTimeCompact("ListenAndServeHttp: %v", currentPortG)
	}

}

// func startHttpsServer(portA string) {
// 	plByMode("https port: %v", portA)

// 	err := http.ListenAndServeTLS(":"+portA, filepath.Join(basePathG, "server.crt"), filepath.Join(basePathG, "server.key"), nil)
// 	if err != nil {
// 		plByMode("ListenAndServeHttps: %v", err.Error())
// 	} else {
// 		plByMode("ListenAndServeHttps: %v", portA)
// 	}

// }

func doCmd(strA string) (string, error) {
	var out bytes.Buffer

	cmd := exec.Command("cmd", "/c", strA)

	cmd.Stdout = &out
	errT := cmd.Run()
	if errT != nil {
		return "", errT
	}

	rStrT := tk.Trim(out.String())

	return rStrT, nil
}

func onceWork() {
	pathT := filepath.Join(basePathG, "once.json")

	fcT := tk.LoadStringFromFile(pathT)

	if tk.IsErrorString(fcT) {
		tk.LogWithTimeCompact("no once work found")

		return
	}

	objT, errT := tk.FromJSON(fcT)

	if errT != nil {
		tk.LogWithTimeCompact("invalid JSON: %v", errT)
		return
	}

	// tk.LogWithTimeCompact("JSON object: %T", objT)

	aryT, ok := objT.([]interface{})

	if !ok {
		tk.LogWithTimeCompact("invalid JSON: %v", "")
		return
	}

	for i, vo := range aryT {
		v, ok := vo.(map[string]interface{})
		if !ok {
			tk.LogWithTimeCompact("invalid JSON item: %#v", vo)
			continue
		}

		cmdStrT, ok := v["Cmd"].(string)
		if !ok {
			tk.LogWithTimeCompact("invalid JSON item: %#v", vo)
			continue
		}

		_, errT = doCmd(cmdStrT)

		if errT != nil {
			tk.LogWithTimeCompact("failed to run once task cmd([%v] %v): %v", i, cmdStrT, errT)
			continue
		}

		tk.LogWithTimeCompact("Running once task [%v] (%v) completed.", i, cmdStrT)

	}

}

type RepeatTask struct {
	Name   string
	Cmd    string
	Start  string
	Period string
}

type RepeatTasks []RepeatTask

var repeatWorksG RepeatTasks = nil

var lastTimesG []time.Time = nil

func doTask(idxA int, taskA RepeatTask) {

	if lastTimesG == nil {
		return
	}

	// default "2006-01-02 15:04:05"
	timeStartT, errT := tk.StrToTimeByFormat(taskA.Start, "")

	if errT != nil {
		tk.LogWithTimeCompact("invalid start time(%v) for %v: %v", taskA.Start, taskA.Name, errT)

		return
	}

	for {
		timeNowT := time.Now()

		diffT := timeNowT.Sub(timeStartT) // / time.Second

		if diffT < 0 {
			tk.SleepSeconds(60)
			continue
		}

		if lastTimesG[idxA].IsZero() {
			lastTimesG[idxA] = timeStartT
		}

		diffT = timeNowT.Sub(lastTimesG[idxA])

		// tk.LogWithTimeCompact("[%v] %v: start time: %v, last time: %v, period: %v, diff: %v", idxA, taskA.Name, taskA.Start, lastTimesG[idxA], taskA.Period, diffT)

		runFlagT := false

		var basePeriodT time.Duration = 0

		switch taskA.Period {
		case "", "minute":

			basePeriodT = time.Second * 60

		case "hour":

			basePeriodT = time.Second * 60 * 60

		case "day":

			basePeriodT = time.Second * 60 * 60 * 24

		case "week":

			basePeriodT = time.Second * 60 * 60 * 24 * 7

		default:

			c, errT := tk.StrToInt(taskA.Period)

			if (errT != nil) || (c <= 0) {
				tk.LogWithTimeCompact("invalid period(%v) for [%v] %v", taskA.Period, idxA, taskA.Name)

				return

			}

			basePeriodT = time.Duration(c) * time.Second * 60

		}

		// taskA.Last = timeNowT

		if basePeriodT <= 0 {
			tk.LogWithTimeCompact("invalid period(%v) for [%v] %v", taskA.Period, idxA, taskA.Name)

			return
		}

		diffInT := diffT / basePeriodT

		if diffInT >= 1 {
			lastTimesG[idxA] = lastTimesG[idxA].Add(diffInT * basePeriodT)

			runFlagT = true
		}

		if runFlagT {
			_, errT = doCmd(taskA.Cmd)

			if errT != nil {
				tk.LogWithTimeCompact("failed to run repeat task cmd([%v] %v): %v", idxA, taskA.Cmd, errT)
			} else {
				tk.LogWithTimeCompact("Running repeat task [%v] (%v) completed.", idxA, taskA.Cmd)
			}

		} else {
			// tk.LogWithTimeCompact("skip diffInT: %d", diffInT)
		}

		tk.SleepSeconds(60)
	}

}

func repeatWork() {
	pathT := filepath.Join(basePathG, "repeat.json")

	fcT := tk.LoadStringFromFile(pathT)

	if tk.IsErrorString(fcT) {
		tk.LogWithTimeCompact("no repeat work found")

		return
	}

	errT := json.Unmarshal([]byte(fcT), &repeatWorksG)

	if errT != nil {
		tk.LogWithTimeCompact("invalid JSON: %v", errT)

		return
	}

	if repeatWorksG == nil {
		return
	}

	lastTimesG = make([]time.Time, len(repeatWorksG))

	for i, v := range repeatWorksG {
		// tk.LogWithTimeCompact("task %v: %v", i, &v)
		go doTask(i, v)
	}
}

func Svc() {
	tk.SetLogFile(filepath.Join(basePathG, serviceNameG+".log"))

	defer func() {
		if v := recover(); v != nil {
			tk.LogWithTimeCompact("panic in svc %v", v)
		}
	}()

	if runModeG != "cmd" {
		runModeG = "service"
	}

	plByMode("runModeG: %v", runModeG)

	tk.DebugModeG = true

	tk.LogWithTimeCompact("------------------")
	tk.LogWithTimeCompact("%v V%v", serviceNameG, versionG)
	tk.LogWithTimeCompact("os: %v, basePathG: %v, configFileNameG: %v", runtime.GOOS, basePathG, configFileNameG)

	if tk.GetOSName() == "windows" {
		plByMode("Windows mode")
		currentOSG = "win"
		if tk.Trim(basePathG) == "" {
			basePathG = "c:\\" + serviceNameG
		}
		configFileNameG = serviceNameG + "win.cfg"
	} else {
		plByMode("Linux mode")
		currentOSG = "linux"
		if tk.Trim(basePathG) == "" {
			basePathG = "/" + serviceNameG
		}
		configFileNameG = serviceNameG + "linux.cfg"
	}

	if !tk.IfFileExists(basePathG) {
		os.MkdirAll(basePathG, 0777)
	}

	tk.SetLogFile(filepath.Join(basePathG, serviceNameG+".log"))

	// currentPortG := "7498"

	cfgFileNameT := filepath.Join(basePathG, configFileNameG)
	if tk.IfFileExists(cfgFileNameT) {
		plByMode("Process config file: %v", cfgFileNameT)
		fileContentT := tk.LoadSimpleMapFromFile(cfgFileNameT)

		if fileContentT != nil {
			currentPortG = fileContentT["port"]
			basePathG = fileContentT["crmBasePath"]
		}
	}

	plByMode("currentPortG: %v, basePathG: %v", currentPortG, basePathG)

	tk.LogWithTimeCompact("currentPortG: %v, basePathG: %v", currentPortG, basePathG)

	tk.LogWithTimeCompact("Service started.")
	tk.LogWithTimeCompact("Using config file: %v", cfgFileNameT)

	// if testPortG > 0 {
	// 	currentPortG = tk.IntToStr(testPortG)
	// 	tk.Pl("currentPortG changed to: %v", currentPortG)
	// }

	go onceWork()

	go repeatWork()

	go startHttpServer(currentPortG)

	// go startHttpsServer(tk.IntToStr(tk.StrToIntWithDefaultValue(currentPortG, 7476) + 1))
}

var exitG = make(chan struct{})

func doWork() {

	go Svc()

	for {
		select {
		case <-exitG:
			os.Exit(0)
			return
		}
	}
}

func runCmd(cmdLineA []string) {
	cmdT := ""

	for _, v := range cmdLineA {
		if !strings.HasPrefix(v, "-") {
			cmdT = v
			break
		}
	}

	// if cmdT == "" {
	// 	fmt.Println("empty command")
	// 	return
	// }

	var errT error

	basePathG = tk.GetSwitchWithDefaultValue(cmdLineA, "-base=", basePathG)

	tk.EnsureMakeDirs(basePathG)

	if !tk.IfFileExists(basePathG) {
		tk.Pl("base path not exists: %v, use current directory instead", basePathG)
		basePathG, errT = filepath.Abs(".")

		if errT != nil {
			tk.Pl("failed to analyze base path")
			return
		}
		// return
	}

	if !tk.IsDirectory(basePathG) {
		tk.Pl("base path not exists: %v", basePathG)
		return
	}

	// tk.Pl("base path: %v", basePathG)

	// testPortG = tk.GetSwitchWithDefaultIntValue(cmdLineA, "-port=", 0)
	// if testPortG > 0 {
	// 	tk.Pl("test port: %v", testPortG)
	// }

	switch cmdT {
	case "version":
		tk.Pl(serviceNameG+" V%v", versionG)
		break
	case "go": // run in cmd mode
		runModeG = "cmd"
		doWork()
		break
	case "test":
		{

		}
		break
	case "", "run":
		s := initSvc()

		if s == nil {
			tk.LogWithTimeCompact("Failed to init service")
			break
		}

		errT = (*s).Run()
		if errT != nil {
			tk.LogWithTimeCompact("Service \"%s\" failed to run: %v.", (*s).String(), errT)
		}
		break
	case "installonly":
		s := initSvc()

		if s == nil {
			tk.Pl("Failed to install")
			break
		}

		errT = (*s).Install()
		if errT != nil {
			tk.Pl("Failed to install: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" installed.", (*s).String())

	case "install":
		s := initSvc()

		if s == nil {
			tk.Pl("Failed to install")
			break
		}

		tk.Pl("Installing service \"%v\"...", (*s).String())

		errT = (*s).Install()
		if errT != nil {
			tk.Pl("Failed to install: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" installed.", (*s).String())

		tk.Pl("Starting service \"%v\"...", (*s).String())

		errT = (*s).Start()
		if errT != nil {
			tk.Pl("Failed to start: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" started.", (*s).String())
	case "uninstall":
		s := initSvc()

		if s == nil {
			tk.Pl("Failed to install")
			break
		}

		errT = (*s).Stop()
		if errT != nil {
			tk.Pl("Failed to stop: %s", errT)
		} else {
			tk.Pl("Service \"%s\" stopped.", (*s).String())
		}

		errT = (*s).Uninstall()
		if errT != nil {
			tk.Pl("Failed to remove: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" removed.", (*s).String())
		break
	case "reinstall":
		s := initSvc()

		if s == nil {
			tk.Pl("Failed to install")
			break
		}

		errT = (*s).Stop()
		if errT != nil {
			tk.Pl("Failed to stop: %s", errT)
		} else {
			tk.Pl("Service \"%s\" stopped.", (*s).String())
		}

		errT = (*s).Uninstall()
		if errT != nil {
			tk.Pl("Failed to remove: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" removed.", (*s).String())

		errT = (*s).Install()
		if errT != nil {
			tk.Pl("Failed to install: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" installed.", (*s).String())

		errT = (*s).Start()
		if errT != nil {
			tk.Pl("Failed to start: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" started.", (*s).String())
	case "restart":
		s := initSvc()

		if s == nil {
			tk.Pl("Failed to install")
			break
		}

		errT = (*s).Stop()
		if errT != nil {
			tk.Pl("Failed to stop: %s", errT)
		} else {
			tk.Pl("Service \"%s\" stopped.", (*s).String())
		}

		errT = (*s).Start()
		if errT != nil {
			tk.Pl("Failed to start: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" restarted.", (*s).String())
	case "start":
		s := initSvc()

		if s == nil {
			tk.Pl("Failed to install")
			break
		}

		errT = (*s).Start()
		if errT != nil {
			tk.Pl("Failed to start: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" started.", (*s).String())
		break
	case "stop":
		s := initSvc()

		if s == nil {
			tk.Pl("Failed to install")
			break
		}

		errT = (*s).Stop()
		if errT != nil {
			tk.Pl("Failed to stop: %v", errT)
			return
		}

		tk.Pl("Service \"%s\" stopped.", (*s).String())
		break
	default:
		tk.Pl("unknown command")
		break
	}

}

func makeBasePath() error {
	if strings.HasPrefix(runtime.GOOS, "win") {
		basePathG = "c:\\" + serviceNameG
	} else {
		basePathG = "/" + serviceNameG
	}

	errT := tk.EnsureMakeDirsE(basePathG)

	if (errT != nil) || (!tk.IfFileExists(basePathG)) {
		basePathG = tk.GetApplicationPath()
	}

	filePathT := filepath.Join(basePathG, "testBasePath.txt")

	errT = tk.SaveStringToFileE("test", filePathT)

	if (errT != nil) || (!tk.IfFileExists(filePathT)) {
		return tk.Errf("init base path failed")
	}

	return nil
}

func main() {

	var errT error

	errT = makeBasePath()

	if errT != nil {
		tk.Pl("failed to init base path: %v", errT)

		return
	}

	if len(os.Args) < 2 {
		tk.Pl("%v V%v is in service(server) mode. Running the application without any arguments will cause it in service mode.\n", serviceNameG, versionG)
		serviceModeG = true

		s := initSvc()

		if s == nil {
			tk.LogWithTimeCompact("Failed to init service")
			return
		}

		err := (*s).Run()
		if err != nil {
			tk.LogWithTimeCompact("Service \"%s\" failed to run.", (*s).String())
		}

		return
	}

	if tk.GetOSName() == "windows" {
		plByMode("Windows mode")
		currentOSG = "win"
		configFileNameG = serviceNameG + "win.cfg"
	} else {
		plByMode("Linux mode")
		currentOSG = "linux"
		configFileNameG = serviceNameG + "linux.cfg"
	}

	if !tk.IfFileExists(basePathG) {
		os.MkdirAll(basePathG, 0777)
	}

	tk.SetLogFile(filepath.Join(basePathG, serviceNameG+".log"))

	// currentPortG := "7498"

	cfgFileNameT := filepath.Join(basePathG, configFileNameG)
	if tk.IfFileExists(cfgFileNameT) {
		plByMode("Process config file: %v", cfgFileNameT)
		fileContentT := tk.LoadSimpleMapFromFile(cfgFileNameT)

		if fileContentT != nil {
			currentPortG = fileContentT["port"]
			basePathG = fileContentT["basePath"]
		}
	}

	plByMode("currentPortG: %v, basePathG: %v", currentPortG, basePathG)

	runCmd(os.Args[1:])

}
