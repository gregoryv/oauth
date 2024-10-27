#!/bin/bash
mkdir -p /var/local/servant
chown gregory /var/local/servant

systemctl stop servant
cp servant /home/gregory/bin/servant
cp servant.service /lib/systemd/system/servant.service
systemctl enable servant
systemctl daemon-reload
systemctl start servant

