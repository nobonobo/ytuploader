# ytuploader

## usage

1. access google cloud console and enable YouTuveDataAPIv3.
2. generate oauth client secret for web application and save "client_secrets.json" to local folder.
3. go get github.com/porjo/youtubeuploader
4. go run github.com/porjo/youtubeuploader -filename=sample.mp4
5. automaticaly open browser and authorize your YouTube account.
6. make
7. scp {youtubeuploader,ytuploader,client_secrets.json,request.token} pi@hostname.local:/home/pi/

## on RaspberryPi

```
sudo apt install samba
```

create "/etc/systemd/system/ytuploader.service"
```
[Unit]
Description=YouTube Auto Uploader

[Service]
WorkingDirectory=/home/pi
ExecStart=/home/pi/ytuploader -src=/home/samba
Restart=always

[Install]
WantedBy=multi-user.target
```

```
sudo systemctl daemon-reload
sudo systemctl start ytuploader
sudo systemctl enable ytuploader
```
