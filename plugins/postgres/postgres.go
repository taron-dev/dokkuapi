package postgres

import (
	"errors"
	"fmt"
	common "github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// SupportedVersions represents allowed postgres versions
var SupportedVersions []string = []string{"11.6"}

var postgresRoot = os.Getenv("POSTGRES_ROOT")

const postgres = "postgres"

// Create create postgres service
func Create(serviceName string, serviceVersion string) (int, string, error) {

	if !isValidServiceName(serviceName) {
		log.ErrorLogger.Println("Invalid service name ", serviceName)
		return 422, "Invalid service name " + serviceName, errors.New("Invalid service name")
	}

	// verify if postgres image exists
	if !imageExists(postgres, serviceVersion) {
		pullImage(postgres, serviceVersion)
	}

	serviceRoot := fmt.Sprintf("%v/%v", postgresRoot, serviceName)
	os.MkdirAll(serviceRoot, 0755)
	os.MkdirAll(serviceRoot+"/data", 0755)
	// LINKS file
	file, err := os.Create(serviceRoot + "/LINKS")
	if err != nil {
		log.ErrorLogger.Println("Cant create LINKS file, ", err)
		// rollback dirs? TODO
		return 500, "LINKS file already exists", err
	}
	defer file.Close()

	//password
	newPwd, err := generatePassword()
	if err != nil {
		log.ErrorLogger.Println(err)
		return 500, "Can't generate pwd", err
	}
	err = ioutil.WriteFile(serviceRoot+"/PASSWORD", newPwd, 0640)
	if err != nil {
		log.ErrorLogger.Println(err)
		return 500, "Can't create /PASSWORD file ", err
	}

	//ENV
	envFile, err := os.Create(serviceRoot + "/ENV")
	if err != nil {
		log.ErrorLogger.Println("Cant create ENV file, ", err)
		// rollback dirs? TODO
		return 500, "ENV file already exists", err
	}
	defer envFile.Close()

	//DATABASE_NAME
	err = ioutil.WriteFile(serviceRoot+"/DATABASE_NAME", []byte(serviceName), 0640)
	if err != nil {
		log.ErrorLogger.Println(err)
		return 500, "Can't create /DATABASE_NAME file ", err
	}

	//TODO service name pwd etc
	// TODO databaseName
	databaseName := serviceName
	log.GeneralLogger.Println("Database name", databaseName)
	serviceName = "dokku.postgres." + serviceName
	log.GeneralLogger.Println("New service name", serviceName)

	err = writeID(serviceRoot, serviceName, serviceVersion, string(newPwd))
	if err != nil {
		log.ErrorLogger.Println(err)
		return 500, "Can't create /ID file ", err
	}

	log.GeneralLogger.Println("Waiting for container to be ready")
	ok := dockerLink(serviceName)
	if !ok {
		log.ErrorLogger.Println("Docker can't link service")
		return 500, "Docker can't link service", errors.New("Docker can't link service")
	}

	log.GeneralLogger.Println("Creating container database")
	ok = createDatabase(serviceName, databaseName)
	if !ok {
		log.ErrorLogger.Println("Docker can't create database container")
		return 500, "Docker can't create database container", errors.New("Docker can't create database container")
	}

	log.GeneralLogger.Println("Securing connection to database")

	ok = stopService(serviceName)
	if !ok {
		log.ErrorLogger.Println("Docker can't stop service")
		return 500, "Docker can't stop service", errors.New("Docker can't stop service")
	}
	ok = secureDbConnection(serviceRoot, postgres, serviceVersion)
	if !ok {
		log.ErrorLogger.Println("Docker can't secure service")
		return 500, "Docker can't secure service", errors.New("Docker can't secure service")
	}

	ok = startPrevious(serviceRoot, serviceName)
	if !ok {
		log.ErrorLogger.Println("Docker can't start previous")
		return 500, "Docker can't start previous", errors.New("Docker can't start previous")
	}

	return 201, "Service created successfully", nil
}

func stopService(serviceName string) bool {
	imageCmd := common.NewShellCmd(strings.Join([]string{"docker", "stop", serviceName}, " "))
	imageCmd.ShowOutput = false
	return imageCmd.Execute()
}

func secureDbConnection(serviceHostRoot string, serviceImage string, serviceVersion string) bool {
	serviceRootBinding := fmt.Sprintf("%v/data:/var/lib/postgresql/data", serviceHostRoot)
	imageLabel := fmt.Sprintf("%v:%v", serviceImage, serviceVersion)
	cmd := []string{
		"docker", "run", "--rm", "-i", "-v", serviceRootBinding, imageLabel,
		"bash", "-s", "</var/lib/dokku/plugins/available/postgres/scripts/enable_ssl.sh"}
	secureCmd := common.NewShellCmd(strings.Join(cmd, " "))
	secureCmd.ShowOutput = false
	return secureCmd.Execute()
}

func startPrevious(serviceRoot string, serviceName string) bool {
	nameFilter := fmt.Sprintf("name=^/%v$", serviceName)
	args := []string{
		"docker", "ps", "-aq", "--no-trunc",
		"--filter", "status=exited",
		"--filter", nameFilter,
		"--format", "{{.ID}}"}
	cmd := common.NewShellCmd(strings.Join(args, " "))
	cmd.ShowOutput = false
	out, err := cmd.Output()
	if err != nil {
		log.ErrorLogger.Println("Get container ID fail:", err.Error())
		return false
	}
	previousId := strings.TrimSpace(string(out))
	log.GeneralLogger.Println("previous ID:", previousId)

	dockerStartPreviousCmd := common.NewShellCmd(strings.Join([]string{
		"docker", "start", previousId}, " "))
	ok := dockerStartPreviousCmd.Execute()
	if !ok {
		log.ErrorLogger.Println("docker start", previousId, "FAIL")
		return false
	}
	log.GeneralLogger.Println("Docker start successfull")
	//TODO service_port_unpause
	return true

}

func isValidServiceName(serviceName string) bool {
	r, _ := regexp.Compile("^[a-z0-9][^:A-Z]*$")
	if r.MatchString(serviceName) {
		return true
	}
	return false
}

func imageExists(imageName string, imageVersion string) bool {
	//docker images -q postgres:11.6
	imageString := strings.Join([]string{imageName, imageVersion}, ":")
	imageCmd := common.NewShellCmd(strings.Join([]string{"docker", "images", "-q", imageString}, " "))
	imageCmd.ShowOutput = false
	return imageCmd.Execute()
}

func pullImage(imageName string, imageVersion string) bool {
	imageString := strings.Join([]string{imageName, imageVersion}, ":")
	imageCmd := common.NewShellCmd(strings.Join([]string{"docker", "pull", imageString}, " "))
	imageCmd.ShowOutput = false
	return imageCmd.Execute()
}

func generatePassword() ([]byte, error) {
	imageCmd := common.NewShellCmd(strings.Join([]string{"openssl", "rand", "-hex", "16"}, " "))
	imageCmd.ShowOutput = false
	out, err := imageCmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func writeID(serviceRoot string, serviceName string, serviceVersion string, newPwd string) error {
	bindVolume := fmt.Sprintf("%v/data:/var/lib/postgresql/data", serviceRoot)
	postgresPwd := fmt.Sprintf("POSTGRES_PASSWORD=%v", newPwd)
	envFile := fmt.Sprintf("--env-file=%v/ENV", serviceRoot)
	containerMetadata := fmt.Sprintf("%v:%v", postgres, serviceVersion)
	cmd := []string{"docker", "run",
		"--name", serviceName,
		"-v", bindVolume,
		"-e", postgresPwd,
		envFile,
		"-d", "--restart", "always",
		"--label", "dokku=service",
		"--label", "dokku.service=postgres", containerMetadata}

	idCmd := common.NewShellCmd(strings.Join(cmd, " "))
	idCmd.ShowOutput = false
	out, err := idCmd.Output()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(serviceRoot+"/ID", out, 0640)
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}

	return nil
}

func dockerLink(serviceName string) bool {
	serviceNameLink := fmt.Sprintf("%v:%v", serviceName, postgres)
	cmd := []string{
		"docker", "run", "--rm", "--link",
		serviceNameLink, os.Getenv("PLUGIN_WAIT_IMAGE"),
		"-p", os.Getenv("PLUGIN_DATASTORE_WAIT_PORT")}
	linkCmd := common.NewShellCmd(strings.Join(cmd, " "))
	linkCmd.ShowOutput = false
	return linkCmd.Execute()
}

func createDatabase(serviceName string, dbName string) bool {
	//docker exec "$SERVICE_NAME" su - postgres -c "createdb -E utf8 $DATABASE_NAME"
	//createDb := fmt.Sprintf("createdb -E utf8 %v\"", dbName)
	cmdString := []string{
		"docker", "exec",
		"-u", "postgres",
		serviceName,
		"createdb", "-E", "utf8", dbName}
	createDbCmd := common.NewShellCmd(strings.Join(cmdString, " "))
	createDbCmd.ShowOutput = false
	return createDbCmd.Execute()
}

func LinkServiceToApp(serviceName string, appName string) bool {
	cmdStrings := []string{
		"dokku", "postgres:link", serviceName, appName}
	createDbCmd := common.NewShellCmd(strings.Join(cmdStrings, " "))
	createDbCmd.ShowOutput = false
	return createDbCmd.Execute()
}
