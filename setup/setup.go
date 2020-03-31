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
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/cavaliercoder/grab"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

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
	return DownloadSoftwareGithub(githubrepo(), target)
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

			/* //run git command on the cli
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
			*/

			//fmt.Println(latest.TagName)
			//https://github.com/src-d/go-git

			_, err := git.PlainClone(target, false, &git.CloneOptions{
				URL:           lurl,
				ReferenceName: plumbing.NewTagReferenceName(latest.TagName),
				SingleBranch:  true,
				Progress:      os.Stdout,
			})
			if err == nil {
				fmt.Println("Succesful clone")
			} else {
				fmt.Println("Error cloning", err)
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
func Exec(donemsg string, cmd *exec.Cmd) {
	var out bytes.Buffer
	cmd.Stdout = &out
	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		log.Println(err, errbuf.String())
	} else {
		log.Println(donemsg, out.String())
	}
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

	preexit := false
	found := false
	terraformsrc := terradl()
	martinipfwdsrc := pfwddl()
	if runtime.GOOS == "linux" {
		relfile := "/etc/lsb-release"
		if _, err := os.Stat(relfile); err == nil {
			file, err := os.Open(relfile)
			if err != nil {
				fmt.Println("Skipped reading release file because of error", err)
			} else {
				fscanner := bufio.NewScanner(file)

				re := regexp.MustCompile(`Ubuntu`)
				for fscanner.Scan() && !found {
					if re.MatchString(fscanner.Text()) {
						found = true
					}
				}
				if found {
					if Confirm("First run? (Do you want to install prereq on an ubuntu system) (y/n)", scanner) {
						if Confirm("Ubuntu has been found, do you want me to run apt-get install -y apache2 mysql-server mysql-client php php-xml composer zip unzip php-mysql? (y)", scanner) {
							preexit = true
							Exec("Updating", exec.Command("apt-get", "update"))
							aptpackages := []string{"apache2", "mysql-server", "php", "php-xml", "composer", "zip", "unzip", "php-mysql"}

							for _, p := range aptpackages {
								fmt.Println("Installing package", p)
								cmd := exec.Command("apt-get", "install", "-y", p)
								Exec("Installation Done", cmd)
							}

							Exec("Enabling rewrite", exec.Command("a2enmod", "rewrite"))
							Exec("Enabling .httpaccess override", exec.Command("sed", "-i", "/<Directory \\/var\\/www\\/>/,/<\\/Directory>/ s/AllowOverride None/AllowOverride All/", "/etc/apache2/apache2.conf"))
							Exec("Restarting apache", exec.Command("/etc/init.d/apache2", "restart"))

						}
						if _, err := os.Stat("/usr/bin/terraform"); os.IsNotExist(err) {
							if Confirm("Do you want me to install terraform? (y)", scanner) {
								_, err := grab.Get("/tmp/terraform.zip", terraformsrc)
								if err != nil {
									log.Fatal(err)
								} else {
									Exec("Unzip done", exec.Command("unzip", "/tmp/terraform.zip", "-d", "/usr/bin"))
									Exec("Chmod +x done", exec.Command("chmod", "+x", "/usr/bin/terraform"))
								}
							}
						} else {
							fmt.Println("Terraform cmd exists")
						}
						if _, err := os.Stat("/usr/bin/martini-pfwd"); os.IsNotExist(err) {
							if Confirm("Do you want me to install martini-pfwd (y)", scanner) {
								_, err := grab.Get("/tmp/martini-pfwd.zip", martinipfwdsrc)
								if err != nil {
									log.Fatal(err)
								} else {
									Exec("Unzip done", exec.Command("unzip", "/tmp/martini-pfwd.zip", "-d", "/usr/bin"))
									Exec("Chmod +x done", exec.Command("chmod", "+x", "/usr/bin/martini-pfwd"))

								}
							}
						} else {
							fmt.Println("Martini-pfwd cmd exists")
						}

					} else {
						fmt.Println("skipping prereq setup")
					}

				}
			}
		}
	}
	mysqlq := `
mysql -u root -p

#MySQL commands:
CREATE DATABASE martini; 
CREATE USER 'martinidbo'@'localhost' IDENTIFIED WITH mysql_native_password BY 'mypasswordthatissupersecret'; 
GRANT ALL ON martini.* TO 'martinidbo'@'localhost'; 
GRANT USAGE ON *.* TO 'martinidbo'@'localhost' WITH MAX_QUERIES_PER_HOUR 0;
`
	if !found {
		fmt.Println("This system is not detected as ubuntu. That doesn't mean it can't work but the prereq setup did not run.")
		fmt.Println("Make sure you manually install all prereq")
		fmt.Println("packages or equal :", "apache2", "mysql-server", "php7", "php-xml", "composer", "zip", "unzip", "php-mysql")
		fmt.Println("enable mod rewrite")
		fmt.Println("enable override for .httpaccess")
		fmt.Println("restart apache")
		fmt.Println("install terraform", terraformsrc)
		fmt.Println("install martini-pfwd ", martinipfwdsrc)
		fmt.Println("Make sure to setup the mysql db and create a database e.g:")
		fmt.Println(mysqlq)

	}
	if preexit {
		fmt.Println("Apt-get has been ran, please make sure that you have setup the mysql db")
		fmt.Println(mysqlq)
		return nil
	}
	fmt.Printf("Where should I install Martini Web [%s]: ", varwww)
	installdir, _ := scanner.ReadString('\n')
	installdir = strings.TrimSpace(installdir)
	if installdir == "" {
		installdir = varwww
	}

	//if Confirm(fmt.Sprintf("Do you want me to download the latest version to %s. It does require the git command to be installed (type y) : ", installdir), scanner) {
	if Confirm(fmt.Sprintf("Do you want me to download the latest version to %s.(type y) : ", installdir), scanner) {
		DownloadSoftware(installdir)
		Exec("Running composer", exec.Command("composer", "--working-dir=/var/www/html", "install"))
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
