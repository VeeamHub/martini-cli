package setup

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/tdewin/martini-cli/core"
	"golang.org/x/crypto/ssh/terminal"
)

func CreateRestToken(dbhost string, dbname string, dblogin string, dbbytePassword []byte, token string, id int64) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dblogin, string(dbbytePassword), dbhost, dbname))
	defer db.Close()

	if err == nil {

		validuntil := time.Now().Add(time.Hour * 1)

		q := fmt.Sprintf("INSERT INTO martini_token (token,validuntil,userid) VALUES ('%s',%d,%d)", token, validuntil.Unix(), id)
		stmtCreate, err := db.Prepare(q)

		if err != nil {
			log.Printf("Problem creating.. %s", q)
			return err
		}
		defer stmtCreate.Close()

		res, err := stmtCreate.Exec()
		if err == nil {
			fmt.Println("Should be ok -> ", res)
		} else {
			log.Printf("Problem creating..")
			return err
		}

	} else {
		log.Printf("Could not connect")
	}
	return err
}

func CreateUser(dbhost string, dbname string, dblogin string, dbbytePassword []byte, userPasswordByte []byte) (error, int64) {
	var id int64

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dblogin, string(dbbytePassword), dbhost, dbname))
	defer db.Close()

	if err == nil {
		//h := md5.New()
		h := sha512.New()
		h.Write([]byte("martini!"))
		h.Write(userPasswordByte)

		stmtCreate, err := db.Prepare(fmt.Sprintf(`
INSERT INTO martini_user (name,email,hashpassword) VALUES ("admin","","%X") 
		`, h.Sum(nil)))

		if err != nil {
			log.Printf("Problem creating..")
			return err, id
		}
		defer stmtCreate.Close()

		res, err := stmtCreate.Exec()
		if err == nil {
			id, _ = res.LastInsertId()
			fmt.Println("Should be ok, user admin with id created -> ", id)

		} else {
			log.Printf("Problem creating..")
			return err, id
		}

	} else {
		log.Printf("Could not connect")
	}
	return err, id
}
func Setup(dbhost string, dbname string, dblogin string, dbbytePassword []byte) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dblogin, string(dbbytePassword), dbhost, dbname))
	defer db.Close()

	if err == nil {

		for _, nq := range GetCreateStatements() {
			stmtCreate, err := db.Prepare(nq.q)
			if err != nil {
				log.Printf("Problem creating.. table %s", nq.n)
				return err
			}
			defer stmtCreate.Close()

			_, err = stmtCreate.Exec()
			if err == nil {
				fmt.Println("Should be ok -> ", nq.n)
			} else {
				log.Printf("Problem creating.. %s", nq.n)
				return err
			}
		}

	} else {
		log.Printf("Could not connect")
	}
	return err
}

func DownloadSoftware(target string) error {
	return DownloadSoftwareGithub("tdewin/rps", target)
}

func DownloadSoftwareGithub(repo string, target string) error {
	var err error
	check := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)

	resp, err := http.Get(check)
	if err == nil {
		var gs []GithubRelease
		var latest GithubRelease
		var a, _ = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(a, &gs)
		if err == nil {
			for _, gr := range gs {
				if !gr.Prerelease && (gr.Published.After(latest.Published)) {
					latest = gr
				}
			}
			fmt.Println("Latest release:", latest.Zipball, latest.Published)
			lurl := fmt.Sprintf("https://github.com/%s.git", repo)
			fmt.Printf("Run\ngit clone --branch %s %s %s\n", latest.TagName, lurl, target)

			cmd := exec.Command("git", "clone", "--branch", latest.TagName, lurl, target)
			var out bytes.Buffer
			cmd.Stdout = &out
			var errbuf bytes.Buffer
			cmd.Stderr = &errbuf

			err := cmd.Run()
			if err != nil {
				log.Println(err, errbuf.String())
			} else {
				log.Println("Seems succesful", out.String())
			}

		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return err
}

func Confirm(q string, pscanner *bufio.Reader) bool {
	fmt.Print(q)
	c, e := (pscanner).ReadString('\n')
	return (e == nil && strings.TrimSpace(c) == "y")
}

func SetupWizard() error {
	var err error

	fmt.Println("Let's get this party started!")
	scanner := bufio.NewReader(os.Stdin)

	varwww := "/var/www/html"
	slash := "/"
	if runtime.GOOS == "windows" {
		varwww = "C:\\inetpub\\wwwroot"
		slash = "\\"
	}

	fmt.Printf("Where should I install Martini Web [%s]: ", varwww)
	installdir, _ := scanner.ReadString('\n')
	installdir = strings.TrimSpace(installdir)
	if installdir == "" {
		installdir = varwww
	}

	if Confirm(fmt.Sprintf("Do you want me to download the latest version to %s. It does require the git command to be installed (type y) : ", installdir), scanner) {
		DownloadSoftware(installdir)
	}

	fmt.Print("What is the database server [127.0.0.1]: ")
	db, _ := scanner.ReadString('\n')
	db = strings.TrimSpace(db)
	if db == "" {
		db = "127.0.0.1"
	}

	fmt.Print("What is the database name [martini]: ")
	dbname, _ := scanner.ReadString('\n')
	dbname = strings.TrimSpace(dbname)
	if dbname == "" {
		dbname = "martini"
	}

	fmt.Print("What is the database login [martinidbo]: ")
	dblogin, _ := scanner.ReadString('\n')
	dblogin = strings.TrimSpace(dblogin)
	if dblogin == "" {
		dblogin = "martinidbo"
	}

	fmt.Print("What is the database password:")
	dbbytePassword, errp := terminal.ReadPassword(int(syscall.Stdin))
	for errp != nil || len(string(dbbytePassword)) < 3 {
		fmt.Println()
		fmt.Print("Password can not be empty (min 3 char):")
		dbbytePassword, errp = terminal.ReadPassword(int(syscall.Stdin))
	}
	fmt.Println()

	if Confirm(fmt.Sprintf("Will setup database %s on instance %s with username %s. Do you want to continue (type y) : ", dbname, db, dblogin), scanner) {
		err = Setup(db, dbname, dblogin, dbbytePassword)

		if err == nil {
			fmt.Print("Type in the admin password: ")
			userPasswordByte, errp := terminal.ReadPassword(int(syscall.Stdin))
			for errp != nil || len(string(userPasswordByte)) < 3 {
				fmt.Println()
				fmt.Print("Password can not be empty (min 3 char):")
				userPasswordByte, errp = terminal.ReadPassword(int(syscall.Stdin))
			}

			err, id := CreateUser(db, dbname, dblogin, dbbytePassword, userPasswordByte)
			if err != nil {
				log.Println("Problem creating user ", err)
			}
			token := uuid.New().String()
			err = CreateRestToken(db, dbname, dblogin, dbbytePassword, token, id)
			if err != nil {
				log.Println("Problem creating token", err)
			} else {
				log.Println("Set temporary token :", token)

				hdir, _ := homedir.Dir()
				cfile := path.Join(hdir, ".martiniconfig")
				var cc core.ClientConfig
				cc.Token = token
				cc.Server = "https://localhost/api"
				cc.Username = "admin"
				jstext, err := json.Marshal(cc)
				if err == nil {
					err = ioutil.WriteFile(cfile, jstext, os.FileMode(0640))
					if err != nil {
						log.Printf("Unable to save config %v", err)
					}
				} else {
					log.Printf("Unable to save config %v", err)
				}
			}

			php := fmt.Sprintf(`<?php 
function getDBParam() {
	return array(
			"dbinstance" => "%s",
			"dbname" => "%s",
			"dbusername" => "%s",
			"dbpassword" => "%s"
	);
}
?>`, db, dbname, dblogin, string(dbbytePassword))

			configfile := strings.Join([]string{installdir, "config.php"}, slash)
			fmt.Println("Created config file in", configfile)

			err = ioutil.WriteFile(configfile, []byte(php), 0644)
			if err != nil {
				log.Println("Can not create config.php, you can manually create it with content\n", php)
				log.Fatal(err)
			}
		} else {
			fmt.Println("Something went wrong in db setup")
		}
	} else {
		fmt.Println("You did not agree with the settings so I'm turning down the music")
	}

	return err
}
