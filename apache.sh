#!/bin/bash

sudo apt update
sudo apt install -y apache2

sudo a2enmod proxy proxy_http

cat <<EOF | sudo tee /etc/apache2/sites-available/mvc.sdslabs.local.conf
<VirtualHost *:80>
    ServerName sdslib.org
    ServerAdmin blaze@mail.com
    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:8000/
    ProxyPassReverse / http://127.0.0.1:8000/
    TransferLog /var/log/apache2/mvc_access.log
    ErrorLog /var/log/apache2/mvc_error.log
</VirtualHost>
EOF

sudo a2ensite mvc.sdslabs.local.conf

if ! grep -q "127.0.0.1 sdslib.org" /etc/hosts; then
    echo "127.0.0.1    sdslib.org" | sudo tee -a /etc/hosts
fi

sudo a2dissite 000-default.conf

sudo apache2ctl configtest

sudo systemctl restart apache2

sudo systemctl status apache2

