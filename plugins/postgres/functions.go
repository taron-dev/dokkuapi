package postgres

import (
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

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

func getServiceUrl(serviceName string) (string, error) {
	password, err := getServicePassword(serviceName)
	if err != nil {
		return "", err
	}
	hostname := "dokku-postgres-" + serviceName
	dbName, err := getDatabaseName(serviceName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("postgres://postgres:%v@%v:%v/%v", password, hostname, postgresPort, dbName), nil
}

func getServicePassword(serviceName string) (string, error) {
	passwordFilePath := fmt.Sprintf("%v/%v/PASSWORD", postgresRoot, serviceName)
	return readFileContent(passwordFilePath)
}

func getDatabaseName(serviceName string) (string, error) {
	dbNameFilePath := fmt.Sprintf("%v/%v/DATABASE_NAME", postgresRoot, serviceName)
	return readFileContent(dbNameFilePath)

}

func readFileContent(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.ErrorLogger.Println("Can read file:", path)
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func isAlreadyLinked(appName string, serviceLink string) bool {
	envFilePath := fmt.Sprintf("%v/%v/ENV", dokkuRoot, appName)
	contentBytes, _ := readFileContent(envFilePath)
	content := string(contentBytes)
	return strings.Contains(content, serviceLink)
}

func appendToFile(path string, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.ErrorLogger.Println("Can't open file:", path)
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(text + "\n"); err != nil {
		log.ErrorLogger.Println("Can't append text")
		return err
	}
	return nil
}

func addLinkToDockerOptions(phases []string, appName string, linkText string) error {
	for _, phase := range phases {
		dockerOptionPhaseFile := fmt.Sprintf("%v/%v/DOCKER_OPTIONS_%v", dokkuRoot, appName, phase)
		err := appendToFile(dockerOptionPhaseFile, linkText)
		if err != nil {
			log.ErrorLogger.Println("Can't append to file:", dockerOptionPhaseFile)
			return err
		}
	}
	return nil
}
