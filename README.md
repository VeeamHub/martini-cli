# Martini CLI

A CLI for [Project Martini](https://github.com/VeeamHub/martini-web) using [Go](https://golang.org/).

## 📗 Documentation

Download CLI

```bash
apt-get install unzip wget

wget http://dewin.me/martini/martini-cli.zip
unzip martini-cli.zip -d /usr/bin
chmod +x /usr/bin/martini-cli
```

On Ubuntu, run first time to install prereq

```bash
martini-cli setup
```

After you made the database (cli will explain you what to do, root password is "" so please change), rerun setup but not first time

```bash
martini-cli setup
```

Once everything is installed, try to connect (--server must be installed if you did not generate certificates)

```bash
martini-cli --server http://localhost/api connect
martini-cli tenant list
```

Please consider (self-signed) certificates
https://www.vultr.com/docs/configure-apache-with-select-signed-tls-ssl-certificate-on-ubuntu-16-04

# Manual prereq install

```bash
apt-get install -y apache2 mysql-server mysql-client php7.2 php7.2-xml composer zip unzip php7.2-mysql
```

```bash
wget https://releases.hashicorp.com/terraform/0.11.14/terraform_0.11.14_linux_amd64.zip
unzip terraform_0.11.14_linux_amd64.zip 
mv terraform /usr/bin
```

```bash
wget http://18.185.97.211:7333/redistr/martini-cli
wget http://18.185.97.211:7333/redistr/martini-pfwd
chmod +x martini*
mv martini* /usr/bin
```

mysql

```bash
mysql -u root -p
```

SQL commands:

```sql
CREATE DATABASE martini; 
CREATE USER 'martinidbo'@'localhost' IDENTIFIED BY 'gkGfLhK6Vbg399q2'; 
GRANT ALL ON martini.* TO 'martinidbo'@'localhost'; 
GRANT USAGE ON *.* TO 'martinidbo'@'localhost' WITH MAX_QUERIES_PER_HOUR 0;
FLUSH privileges;
```

enable rewrite mod

```bash
a2enmod rewrite
```

enable override, open it with for example nano

```bash
nano /etc/apache2/apache2.conf
```

```bash
<Directory /var/www/>
        Options Indexes FollowSymLinks
        AllowOverride none
        Require all granted
</Directory>
```

with:

```bash
<Directory /var/www/>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
</Directory>
```

restart service

```bash
/etc/init.d/apache2 restart
```

## ✍ Contributions

We welcome contributions from the community! We encourage you to create [issues](https://github.com/VeeamHub/martini-cli/issues/new/choose) for Bugs & Feature Requests and submit Pull Requests for improving our documentation. For more detailed information, refer to our [Contributing Guide](CONTRIBUTING.md).

## 🤝🏾 License

* [MIT License](LICENSE)

## 🤔 Questions

If you have any questions or something is unclear, please don't hesitate to [create an issue](https://github.com/VeeamHub/martini-cli/issues/new/choose) and let us know!
