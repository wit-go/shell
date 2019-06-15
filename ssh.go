package shell

import "log"
import "fmt"
import "io/ioutil"
import "time"

// import "strings"
// import "path/filepath"
// import "os"
// import "bufio"
// import "os/user"
// import "runtime"

import "golang.org/x/crypto/ssh"
import "github.com/tmc/scp"

var sshHostname string
var sshPort	int
var sshUsername	string
var sshPassword	string
var sshKeyfile	string

func SSHclientSet(hostname string, port int, username string, pass string, keyfile string) {
	sshHostname	= hostname
	sshPort		= port
	sshUsername	= username
	sshPassword	= pass
	sshKeyfile	= keyfile
}

func SSHclientSCP(localfile string, remotefile string) {
	log.Println("shell.SSHclientSCP() START")
	log.Println("shell.SSHclientSCP() \tlocalfile =", localfile)
	log.Println("shell.SSHclientSCP() \tremotefile =", remotefile)
	sess := SSH(sshHostname, sshPort, sshUsername, sshPassword, sshKeyfile)
	err := scp.CopyPath(localfile, remotefile, sess)
	sess.Close()
	log.Println("shell.SSHclientSCP() \tscp.CopyPath() err =", err)
	log.Println("shell.SSHclientSCP() END")
}

func SSHclientRun(cmd string) {
	log.Println("shell.SSHclientRun() START cmd =", cmd)
	sess := SSH(sshHostname, sshPort, sshUsername, sshPassword, sshKeyfile)
	err := sess.Run(cmd)
	sess.Close()
	log.Println("shell.SSHclientRun() END err =", err)
}

func SSH(hostname string, port int, username string, pass string, keyfile string) *ssh.Session {
	// get host public key
	// hostKey := getHostKey(host)
	// log.Println("hostkey =", hostKey)

	publicKey, err := PublicKeyFile(keyfile)
	if (err != nil) {
		log.Println("PublicKeyFile() error =", err)
	}

	// ssh client config
	config := ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
			publicKey,
		},
		// allow any host key to be used (non-prod)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		// verify host public key
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
		// optional host key algo list
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		// optional tcp connect timeout
		Timeout:         5 * time.Second,
	}

	sport := fmt.Sprintf("%d", port)
	// connect
	client, err := ssh.Dial("tcp", hostname+":"+sport, &config)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Close()

	// start session
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	// defer sess.Close()

	return sess
}

func Scp(sess *ssh.Session, localfile string, remotefile string) {
	err := scp.CopyPath(localfile, remotefile, sess)
	log.Println("scp.CopyPath() err =", err)
}

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	log.Println("buffer =", string(buffer))
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

// THIS doesn't work
/*
func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	// 9enFJdMhb8eHN/6qfHSU/jww2Mo=|pcsWQCvAyve9QXBhjL+w/LhkcHU= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBMQx8BJXxD+vk3wyjy7Irzw4FA6xxJvqUP7Hb+Z+ygpOuidYj9G8x6gHEXFUnABn5YirePrWh5tNsk4Rqs48VwU=
	hostKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte("9enFJdMhb8eHN/6qfHSU/jww2Mo=|pcsWQCvAyve9QXBhjL+w/LhkcHU= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBMQx8BJXxD+vk3wyjy7Irzw4FA6xxJvqUP7Hb+Z+ygpOuidYj9G8x6gHEXFUnABn5YirePrWh5tNsk4Rqs48VwU="))
	log.Println("hostkey err =", err)
	log.Println("hostkey =", hostKey)
	if hostKey == nil {
		log.Println("no hostkey found err =", err)
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}
*/
