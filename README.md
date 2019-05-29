# Martini-cli

Download cli
```
apt-get install unzip wget

wget http://dewin.me/martini/martini-cli.zip
unzip martini-cli.zip -d /usr/bin
chmod +x /usr/bin/martini-cli
```

On ubuntu, run first time to install prereq
```
martini-cli setup
```

After you made the database (cli will explain you what to do, root password is "" so please change), rerun setup but not first time
```
martini-cli setup
```

Once everything is install, try to connect (--server must be installed if you did not generate certificates)
```
martini-cli --server http://localhost/api connect
martini-cli tenant list
```

Please consider (selfsigned) certificates
https://www.vultr.com/docs/configure-apache-with-select-signed-tls-ssl-certificate-on-ubuntu-16-04



# Manual prereq install
```
apt-get install -y apache2 mysql-server mysql-client php7.2 php7.2-xml composer zip unzip php7.2-mysql
```

```
wget https://releases.hashicorp.com/terraform/0.11.14/terraform_0.11.14_linux_amd64.zip
unzip terraform_0.11.14_linux_amd64.zip 
mv terraform /usr/bin
```

```
wget http://18.185.97.211:7333/redistr/martini-cli
wget http://18.185.97.211:7333/redistr/martini-pfwd
chmod +x martini*
mv martini* /usr/bin
```

mysql 
```
mysql -u root -p
```

SQL commands:
```
CREATE DATABASE martini; 
CREATE USER 'martinidbo'@'localhost' IDENTIFIED BY 'gkGfLhK6Vbg399q2'; 
GRANT ALL ON martini.* TO 'martinidbo'@'localhost'; 
GRANT USAGE ON *.* TO 'martinidbo'@'localhost' WITH MAX_QUERIES_PER_HOUR 0;
FLUSH privileges;
```

enable rewrite mod
```
a2enmod rewrite
```

enable override, open it with for example nano
```
nano /etc/apache2/apache2.conf
```
```
<Directory /var/www/>
        Options Indexes FollowSymLinks
        AllowOverride none
        Require all granted
</Directory>
```
with:
```
<Directory /var/www/>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
</Directory>
```

restart service
```
/etc/init.d/apache2 restart
```