package shell

import "bufio"
import "log"
import "fmt"
import "os"
import "os/user"
import "io/ioutil"
import "path/filepath"
import "strings"
import "time"
import "runtime"

import "golang.org/x/crypto/ssh"
import "github.com/tmc/scp"

func SSH(hostname string, port int, username string, pass string) *ssh.Session {
//	username := "jcarr"
//	pass := "tryme"
	// cmd  := "ps"

	// get host public key
	// hostKey := getHostKey(host)
	// log.Println("hostkey =", hostKey)

	user, _ := user.Current()

	keyfile := user.HomeDir + "/.ssh/id_ed25519"
	if runtime.GOOS == "windows" {
		if Exists("/cygwin") {
			log.Println("On Windows, but running within cygwin")
			keyfile = "/home/wit/.ssh/id_ed25519"
		} else {
			keyfile = user.HomeDir + "\\id_ed25519"
		}
	}

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
