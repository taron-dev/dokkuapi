package ssh

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// AddSSHPublicKey add ssh key to machine
func AddSSHPublicKey(userName string, publicKey string) bool {
	var dokkuRoot = os.Getenv("DOKKU_ROOT")
	pubKeyPath := dokkuRoot + "/.ssh/" + userName + ".pub"
	authorizedKeysPath := dokkuRoot + "/.ssh/authorized_keys"

	//dat kluc do suboru
	err := storeKeyToFile(pubKeyPath, publicKey)
	if err != nil {
		return false
	}
	//Validacia suboru
	if !validateKeyFile(pubKeyPath) {
		return false
	}

	authorizedKeyFile, _ := os.OpenFile(authorizedKeysPath, os.O_RDONLY|os.O_CREATE, 0755)
	authorizedKeyFile.Close()

	//sshcommand
	args := []string{"sshcommand", "acl-add", "dokku", userName, pubKeyPath}
	cmd := common.NewShellCmd(strings.Join(args, " "))
	cmd.ShowOutput = false
	out, err := cmd.Output()
	if err != nil {
		os.Remove(pubKeyPath)
		log.ErrorLogger.Println("Can't add sshkey:", err.Error())
		return false
	}
	log.GeneralLogger.Println(string(out))

	// verify authorized_key file
	authorizedKeyFile, err = os.Open(authorizedKeysPath)
	if err != nil {
		log.ErrorLogger.Println("Add ssh pub key failed after ssh-keygen check:", err.Error())
		return false
	}
	defer authorizedKeyFile.Close()
	scanner := bufio.NewScanner(authorizedKeyFile)
	scanner.Split(bufio.ScanLines)
	tmpKeyPath := dokkuRoot + "/.ssh/" + userName
	for scanner.Scan() {
		err := storeKeyToFile(tmpKeyPath, scanner.Text())
		if err != nil {
			log.ErrorLogger.Println("Can't store TMP ket to:", tmpKeyPath)
		}
		if !validateKeyFile(pubKeyPath) {
			log.ErrorLogger.Println(authorizedKeysPath, "validation failed")
			return false
		}
	}
	os.Remove(tmpKeyPath)
	return true
}

// RemoveSSHPublicKey removes user related ssh public key
func RemoveSSHPublicKey(userName string) bool {
	var dokkuRoot = os.Getenv("DOKKU_ROOT")
	pubKeyPath := dokkuRoot + "/.ssh/" + userName + ".pub"
	os.Remove(pubKeyPath)
	args := []string{"sshcommand", "acl-remove", "dokku", userName}
	cmd := common.NewShellCmd(strings.Join(args, " "))
	cmd.ShowOutput = false
	return cmd.Execute()
}

// UserHasPublicSSHKey check if user has already added public ssh key
func UserHasPublicSSHKey(userName string) (bool, error) {
	out, err := exec.Command("dokku", "ssh-keys:list").CombinedOutput()
	output := string(out)
	if err != nil {
		log.ErrorLogger.Println("Can't get public ssh keys list:", err.Error(), output)
		return false, err
	}
	regex := regexp.MustCompile("NAME=\"" + userName + "\"")
	matches := regex.FindStringSubmatch(output)
	if len(matches) > 0 {
		return true, nil
	}

	return false, nil
}

// IsValidPublicSSHKey validate publicKey
func IsValidPublicSSHKey(userName string, publicKey string) (bool, error) {
	uuid, err := newUUID()
	if err != nil {
		return false, err
	}
	var dokkuRoot = os.Getenv("DOKKU_ROOT")
	pubKeyPath := dokkuRoot + "/.ssh/" + userName + uuid + ".pub"
	//dat kluc do suboru
	err = storeKeyToFile(pubKeyPath, publicKey)
	if err != nil {
		return false, err
	}
	//Validacia suboru
	if !validateKeyFile(pubKeyPath) {
		return false, nil
	}
	os.Remove(pubKeyPath)
	return true, nil
}

func storeKeyToFile(path string, publicKey string) error {
	err := ioutil.WriteFile(path, []byte(publicKey), 0755)
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func validateKeyFile(path string) bool {
	args := []string{"ssh-keygen", "-lf", path}
	cmd := common.NewShellCmd(strings.Join(args, " "))
	cmd.ShowOutput = false
	if isValid := cmd.Execute(); !isValid {
		log.ErrorLogger.Println("Public key is invalid")
		//remove created file
		os.Remove(path)
		return false
	}
	return true
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x", uuid), nil
}
