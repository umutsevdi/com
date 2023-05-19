# umutsevdi.com

This directory contains source code of my server-side rendered web application.

```
webserver/
├── app/
│   ├── app.config
│   ├── client/
│   ├── content/
│   │   ├── favicon.ico
│   │   ├── index.html
│   │   ├── not-found.html
│   │   ├── resume.html
│   │   ├── robots.txt
│   │   └── static/
│   │       ├── css/
│   │       ├── fonts/
│   │       ├── img/
│   │       ├── js/
│   │       └── other/
│   ├── go.mod
│   ├── main.go
│   ├── min.sh
│   ├── router.go
│   ├── syslog/
│   └── util/
├── LICENSE
└── README.md
```

## Installing
Install go, and compile the source code
```
cd webserver
wget https://go.dev/dl/go-version
tar -xvf go-version.tar.gz
cd app
../go/bin/go build .
./server
```
Configure the `app.config` file.

## Configuration
To run this program you need a file called `app.config` in the `app/` directory
An `app.config` file consists of keys and values delimited with equal sign.

It has the following parameters:
```c
- site*        /* the domain of the website */
- ip           /* Ip address of the website, default: localhost */
- port         /* Port of the program, default: 8080 */
- token*       /* Github API token to fetch the data */
- user*        /* Github user */
- cache        /* Whether filesystem-caching is enabled or not */
- cache-time   /* How often the caching of GitHub and filesystem should be refreshed*/
```
