package main

import (
	fptr10 "clientrabbit/fptr"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var clearLogsProgramm = flag.Bool("clearlogs", true, "очистить логи программы")
var commandForAction = flag.String("command", "closeShift", "команды: closeShift, reportX, openShift")
var kassirName = flag.String("kassir", "", "имя кассира")
var debugpr = flag.Bool("debug", false, "дебажим программу")
var emulation = flag.Bool("emul", false, "эмуляция")

var LOGSDIR = "./logs/"
var filelogmap map[string]*os.File
var logsmap map[string]*log.Logger

const LOGINFO = "info"
const LOGINFO_WITHSTD = "info_std"
const LOGERROR = "error"
const LOGSKIP_LINES = "skip_line"
const LOGOTHER = "other"
const LOG_PREFIX = "TASKS"
const Version_of_program = "2024_01_21_01"

const FILE_NAME_CONNECTION = "connection.txt"
const FILE_NAME_KASSIRNAME = "kassir.txt"

type TOperator struct {
	Name  string `json:"name"`
	Vatin string `json:"vatin,omitempty"`
}

type TRepotyAtolKKT struct {
	Type     string    `json:"type"`
	Operator TOperator `json:"operator"`
}

func main() {
	var err error
	var descrError string
	var reportAtol TRepotyAtolKKT
	runDescription := "программа версии " + Version_of_program + " закрытия кассовой смены по расписанию"
	fmt.Println(runDescription)
	defer fmt.Println("программа версии " + Version_of_program + " закрытия кассовой смены по расписанию остановлена")

	fmt.Println("парсинг параметров запуска программы")
	flag.Parse()
	fmt.Println("emulation", *emulation)
	fmt.Println("debugpr", *debugpr)
	clearLogsDescr := fmt.Sprintf("Очистить логи программы %v", *clearLogsProgramm)
	fmt.Println(clearLogsDescr)
	fmt.Println("инициализация лог файлов программы")
	if foundedLogDir, _ := doesFileExist(LOGSDIR); !foundedLogDir {
		os.Mkdir(LOGSDIR, 0777)
	}
	filelogmap, logsmap, descrError, err = initializationLogs(*clearLogsProgramm, LOGINFO, LOGERROR, LOGSKIP_LINES, LOGOTHER)
	defer func() {
		fmt.Println("закрытие дескрипторов лог файлов программы")
		for _, v := range filelogmap {
			if v != nil {
				v.Close()
			}
		}
	}()
	if err != nil {
		descrMistake := fmt.Sprintf("ошибка инициализации лог файлов %v", descrError)
		fmt.Fprint(os.Stderr, descrMistake)
		log.Panic(descrMistake)
	}
	fmt.Println("лог файлы инициализированы в папке " + LOGSDIR)
	multwriterLocLoc := io.MultiWriter(logsmap[LOGINFO].Writer(), os.Stdout)
	logsmap[LOGINFO_WITHSTD] = log.New(multwriterLocLoc, LOG_PREFIX+"_"+strings.ToUpper(LOGINFO)+" ", log.LstdFlags)
	logsmap[LOGINFO].Println(runDescription)
	logsmap[LOGINFO].Println(clearLogsDescr)
	defer logsmap[LOGINFO].Println("программа версии " + Version_of_program + " закрытия кассовой смены по расписанию остановлена")

	logsmap[LOGINFO_WITHSTD].Println("инициализация драйвера атол")
	fptr, err := fptr10.NewSafe()
	if err != nil {
		descrError := fmt.Sprintf("Ошибка (%v) инициализации драйвера ККТ атол", err)
		logsmap[LOGERROR].Println(descrError)
		log.Panic(descrError)
	}
	defer fptr.Destroy()
	fmt.Println(fptr.Version())
	//сединение с кассой
	logsmap[LOGINFO_WITHSTD].Println("соединение с кассой")
	kassirNameValue := *kassirName
	if kassirNameValue == "" {
		kassirNameValue, err = getCurrentKassirName(".")
		if err != nil {
			descrError := fmt.Sprintf("ошибка (%v) чтения имени кассира", err)
			logsmap[LOGERROR].Println(descrError)
			kassirNameValue = "Кассир ИО"
		}
	}
	if kassirNameValue == "" {
		logsmap[LOGERROR].Println("имя кассира не определено")
		kassirNameValue = "Кассир ИО"
	}
	comPort, err := getCurrentPortOfKass(".")
	if err != nil {
		desrErr := fmt.Sprintf("ошибка (%v) чтения параметра com порт соединения с кассой", err)
		logsmap[LOGINFO].Println(desrErr)
		comPort = ""
	}
	if !connectWithKassa(fptr, comPort) && !(*emulation) {
		descrErr := fmt.Sprintf("ошибка соединения с кассовым аппаратом на ком порт %v", comPort)
		logsmap[LOGINFO].Println("соедиенние с кассой не установлено")
		logsmap[LOGERROR].Println(descrErr)
		if !*debugpr {
			log.Panic(descrErr)
		}
	} else {
		logsmap[LOGINFO_WITHSTD].Printf("подключение к кассе на порт %v прошло успешно", comPort)
	}
	defer fptr.Close()

	reportAtol.Type = *commandForAction
	reportAtol.Operator.Name = kassirNameValue

	jsonCloseShift, err := json.Marshal(reportAtol)
	if err != nil {
		descrError := fmt.Sprintf("ошибка (%v) формирования json для отправки команды", err)
		logsmap[LOGERROR].Println(descrError)
		log.Panic(descrError)
	}

	resulOfCommand, err := sendComandeAndGetAnswerFromKKT(fptr, string(jsonCloseShift))
	if err != nil {
		errorDescr := fmt.Sprintf("ошибка (%v) закрытия смены на кассе атол", err)
		logsmap[LOGERROR].Println(errorDescr)
		log.Panic(errorDescr)
	}
	if !successCommand(resulOfCommand) {
		logsmap[LOGINFO].Printf("закрытие смены на кассе не прошло. ошибка %v", resulOfCommand)
		log.Panic(resulOfCommand)
	}
	logsmap[LOGINFO].Print("закрытие смены прошлом успешно")
}

func sendComandeAndGetAnswerFromKKT(fptr *fptr10.IFptr, comJson string) (string, error) {
	var err error
	logsmap[LOGINFO].Printf("начало процедуры запуска json задания: %v", comJson)
	//return "", nil
	fptr.SetParam(fptr10.LIBFPTR_PARAM_JSON_DATA, comJson)
	//fptr.ValidateJson()
	if !*emulation {
		err = fptr.ProcessJson()
	}
	if err != nil {
		if !*debugpr {
			desrError := fmt.Sprintf("ошибка (%v) закрытия смены", err)
			logsmap[LOGERROR].Println(desrError)
			//logsmap[LOGINFO].Print("начало процедуры sendComandeAndGetAnswerFromKKT c ошибкой", err)
			return desrError, err
		}
	}
	result := fptr.GetParamString(fptr10.LIBFPTR_PARAM_JSON_DATA)
	logsmap[LOGOTHER].Println("comJson", comJson)
	logsmap[LOGOTHER].Println("result", result)
	return result, nil
}

func doesFileExist(fullFileName string) (found bool, err error) {
	found = false
	if _, err = os.Stat(fullFileName); err == nil {
		// path/to/whatever exists
		found = true
	}
	return
}

func initializationLogs(clearLogs bool, logstrs ...string) (map[string]*os.File, map[string]*log.Logger, string, error) {
	var reserr, err error
	reserr = nil
	filelogmapLoc := make(map[string]*os.File)
	logsmapLoc := make(map[string]*log.Logger)
	descrError := ""
	for _, logstr := range logstrs {
		filenamelogfile := logstr + "logs.txt"
		preflog := LOG_PREFIX + "_" + strings.ToUpper(logstr)
		fullnamelogfile := LOGSDIR + filenamelogfile
		filelogmapLoc[logstr], logsmapLoc[logstr], err = intitLog(fullnamelogfile, preflog, clearLogs)
		if err != nil {
			descrError = fmt.Sprintf("ошибка инициализации лог файла %v с ошибкой %v", fullnamelogfile, err)
			fmt.Fprintln(os.Stderr, descrError)
			reserr = err
			break
		}
	}
	return filelogmapLoc, logsmapLoc, descrError, reserr
}

func successCommand(resulJson string) bool {
	res := true
	indOsh := strings.Index(resulJson, "ошибка")
	indErr := strings.Index(resulJson, "error")
	if indErr != -1 || indOsh != -1 {
		res = false
	}
	return res
}

func connectWithKassa(fptr *fptr10.IFptr, comport string) bool {
	sComPorta := comport
	if !strings.Contains(comport, "COM") {
		sComPorta = "COM" + comport
	}
	fptr.SetSingleSetting(fptr10.LIBFPTR_SETTING_MODEL, strconv.Itoa(fptr10.LIBFPTR_MODEL_ATOL_AUTO))
	if comport != "" {
		fptr.SetSingleSetting(fptr10.LIBFPTR_SETTING_PORT, strconv.Itoa(fptr10.LIBFPTR_PORT_COM))
		fptr.SetSingleSetting(fptr10.LIBFPTR_SETTING_COM_FILE, sComPorta)
	} else {
		fptr.SetSingleSetting(fptr10.LIBFPTR_SETTING_PORT, strconv.Itoa(fptr10.LIBFPTR_PORT_USB))
	}
	fptr.SetSingleSetting(fptr10.LIBFPTR_SETTING_BAUDRATE, strconv.Itoa(fptr10.LIBFPTR_PORT_BR_115200))
	fptr.ApplySingleSettings()
	fptr.Open()
	return fptr.IsOpened()
}

func getCurrentPortOfKass(dirOfJsons string) (string, error) {
	comportb, err := os.ReadFile(dirOfJsons + "\\" + FILE_NAME_CONNECTION)
	if err != nil {
		desrError := fmt.Sprintf("ошибка (%v) октрытия файла с параметрами соедиения кассы", err)
		logsmap[LOGINFO].Println(desrError)
		return desrError, err
	}
	return string(comportb), nil
}

func getCurrentKassirName(dirOfJsons string) (string, error) {
	kassname, err := os.ReadFile(dirOfJsons + "\\" + FILE_NAME_KASSIRNAME)
	if err != nil {
		desrError := fmt.Sprintf("ошибка (%v) октрытия файла с именем кассира", err)
		logsmap[LOGINFO].Println(desrError)
		return desrError, err
	}
	return string(kassname), nil
}

func intitLog(logFile string, pref string, clearLogs bool) (*os.File, *log.Logger, error) {
	flagsTempOpen := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	if clearLogs {
		flagsTempOpen = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(logFile, flagsTempOpen, 0644)
	multwr := io.MultiWriter(f)
	//if pref == LOG_PREFIX+"_INFO" {
	//	multwr = io.MultiWriter(f, os.Stdout)
	//}
	flagsLogs := log.LstdFlags
	if pref == LOG_PREFIX+"_ERROR" {
		multwr = io.MultiWriter(f, os.Stderr)
		flagsLogs = log.LstdFlags | log.Lshortfile
	}
	if err != nil {
		fmt.Println("Не удалось создать лог файл ", logFile, err)
		return nil, nil, err
	}
	loger := log.New(multwr, pref+" ", flagsLogs)
	return f, loger, nil
}
