package postgres

import (
	"errors"
	"fmt"
	"github.com/dokku/dokku/plugins/config"
	log "github.com/ondro2208/dokkuapi/logger"
	"io/ioutil"
	"os"
)

const postgres = "postgres"

var postgresRoot = os.Getenv("POSTGRES_ROOT")
var postgresPort = os.Getenv("POSTGRES_DATASTORE_PORT")
var dokkuRoot = os.Getenv("DOKKU_ROOT")

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

// LinkServiceToApp link service to app
func LinkServiceToApp(serviceName string, appName string) error {
	serviceURL, err := getServiceUrl(serviceName)
	if err != nil {
		return err
	}
	if isAlreadyLinked(appName, serviceURL) {
		return errors.New("Is already linked")
	}
	linksFilePath := fmt.Sprintf("%v/%v/LINKS", postgresRoot, serviceName)
	err = appendToFile(linksFilePath, appName)
	if err != nil {
		return err
	}
	linkText := fmt.Sprintf("--link %v:dokku-postgres-%v", serviceName, serviceName)
	err = addLinkToDockerOptions([]string{"BUILD", "DEPLOY", "RUN"}, appName, linkText)
	if err != nil {
		return err
	}
	return config.SetMany(appName, map[string]string{"DATABASE_URL": serviceURL}, true)
}
