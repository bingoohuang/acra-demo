# Acro Demo

1. create key store for the acra-server: `docker-compose -f docker-infrastructure.yml up acra-keymaker_server`
2. start the database and the acra-server
   - for PostgreSQL: `docker-compose -f docker-pg.yml up`
   - for MySQL: `docker-compose -f docker-my.yml up`
3. insert data and query
   - for PostgreSQL
       - insert data through acra-server `go run demo/demo.go -insert 3`
       - query data through acra-server `go run demo/demo.go -select 3`
       - query data through database directly `go run demo/demo.go -select 3 --dbport 15432`
   - for MySQL
       - insert data through acra-server `go run demo/demo.go -mysql -insert 3`
       - query data through acra-server `go run demo/demo.go -mysql -select 3`
       - query data through database directly `go run demo/demo.go -mysql -select 3 --dbport 15432`

logs for the query data through acra-server:

```sh
$ go run demo/demo.go -mysql -select 3
1	ahoniatzur	3MIhjHos%n8	kyocxhgif@uqbdjx.tmsouz.com
2	bokcimaoj	ZWS3G@3e$-E	oagbiufjxl@icjgsy.dcyygujs.org
3	perbatomiu	VntmIZ7zw19	oqmzku@sqaxszyji.mzpyneh.com
2022/08/10 14:19:22 Select has been successful
```

logs for the query data through database directly:

```sh
$ go run demo/demo.go -select 3 -mysql --dbport 13306
1	ahoniatzur	JSUloQAAAAAAAADwIiIiIpEAAAAAAAAAAKA2AEwAAAEBQAwAAAAQAAAAIAAAAKq4XU5irJ_56G76rARMG_jpe-_be5K_XHoeywPEBicHlhChRh-bHiwCdOVWD4xbbXSxJizMoavY8gxd-QABAUAMAAAAEAAAAAsAAAC_QmFiFdVfZFOwVg6AQHn_VgIQkcSDtF1nNdAJwHNaDTU8K89B9Vg=	JSUlsQAAAAAAAADwIiIiIqEAAAAAAAAAAKA2AEwAAAEBQAwAAAAQAAAAIAAAAK7-S5f0DDEuI602j5P22OdPzEB2uw1LnEjnA31QQjrMW-amN8n53pZ_AXr9FEUGkgh2q4qRCiZ1AZKy9AABAUAMAAAAEAAAABsAAADQnOTCjskvZUx01fBFrigQBlG1SsV-H0LAygfxr8-nCclnzusYWfi6wk171dokh9xmrCB7WYDD
2	bokcimaoj	JSUloQAAAAAAAADwIiIiIpEAAAAAAAAAAKA2AEwAAAEBQAwAAAAQAAAAIAAAAPqKEkI9U5FLOfnLAhxiM4bNXQAAGDj83ZYX9cVT1GmRBJT6zlv_FFsIAzqJH5KISiw_1ibfOkZOijvxZQABAUAMAAAAEAAAAAsAAABdFDEF_hjtg_WRwDG4Q0QtCLzLhXlQUqdPETS9taPrx0mcH0uS_30=	JSUltAAAAAAAAADwIiIiIqQAAAAAAAAAAKA2AEwAAAEBQAwAAAAQAAAAIAAAAB-Ne9xa5gPji1ILQCSCpl7vH5A_7CzvvzH48HRGfr8HOD7O0l8nR3Tjx0zxOR-6BfsrfXeVN1G88Gta7gABAUAMAAAAEAAAAB4AAAAzI754L1FC1IWoL2cY7sBT224poqigkeeaqQwSs9W-04aCqvrKEUMVF1pGn6lx_gGkZgoIOYXzwpHa
3	perbatomiu	JSUloQAAAAAAAADwIiIiIpEAAAAAAAAAAKA2AEwAAAEBQAwAAAAQAAAAIAAAAH2mcZ6I_gjsl2oh-3hOjwQl_9ghYSQ6w3o7hMZ-u-OYjm-Bt-KO0SQTwF2fWKSaogNzVbHD_GcV044o1QABAUAMAAAAEAAAAAsAAAAIb-Es2z5sPszPV64GUIcZkCUEXTNfxNrApxG1kE53BxFkqV-K5IU=	JSUlsgAAAAAAAADwIiIiIqIAAAAAAAAAAKA2AEwAAAEBQAwAAAAQAAAAIAAAAHhTakXdsxIz3hxiNNb04wPDI0ZmUq5X1UNfNMjfnZ4q-6VAZ01jDvEX8H9g4ZIi32CbDBwTvtPxhAqfZwABAUAMAAAAEAAAABwAAADfDHDFdqzSySH2uZybGz4DV4oWbOlYj3BtIMz4rek4I42FOBfRFRywxVwgCATyNV3Mwb6i4b4uMQ==
2022/08/10 14:19:54 Select has been successful
```

## Acra Poison Records Demo

Acra Poison Records Demo illustrates how to use intrusion detection functionality of [Acra data protection suite](https://cossacklabs.com/acra/). For intrusion detection, Acra uses poison records, also known as honey tokens. This demo shows how to setup, configure, and use intrusion detection in Acra.

This project is one of numerous Acra’s example applications. If you are curious about other Acra features, like transparent encryption, SQL firewall, load balancing support – [Acra Example Applications](https://github.com/cossacklabs/acra-engineering-demo/).

## How poison records work?

Poison records are records specifically designed to sit quietly in the database and not be queried by legitimate users under normal circumstances. They look like any other encrypted records, and it’s impossible to distinguish them from “normal data”. Technically speaking, poison records are data (binary or strings, int, or whatever suits your database design), placed in particular tables / columns / cells.

However, poison records will only be included in the outputs of requests from malicious applications that read more data than they should, i.e. using `SELECT *` requests. The sole purpose of these requests is that when an unauthorised leakage occurs, poison records will be present in database response and detected by AcraServer. AcraServer will inform user (system administrator) of untypical behaviour and can block suspicious requests.

Related blog posts and docs:

- [Explain Like I’m Five: Poison Records (Honeypots for Database Tables)](https://hackernoon.com/poison-records-acra-eli5-d78250ef94f)
- [Poison Records In Acra – Database Honeypots For Intrusion Detection](https://www.cossacklabs.com/blog/acra-poison-records.html)
- [Acra docs](https://docs.cossacklabs.com/pages/intrusion-detection/) on intrusion detection

## How to run the example project

### Installation

1. Use docker-compose command to set up and run the whole infrastructure:

    ```bash
    docker-compose -f docker-infrastructure.yml up
    ```

    This will deploy PostgreSQL database and AcraServer in [transparent mode](https://github.com/cossacklabs/acra#integrating-server-side-encryption-using-acraserver-in-transparent-proxy-mode) of operations.

2. Let’s check that those containers are running:

```bash
docker ps -a
```

You should see two containers up and running, and another two in “exited” state (`acra-keymaker` and `acra-poisonrecordmaker`). These containers were used to generate encryption keys for data and poison records themselves. They finished their mission and stopped.

```sh
CONTAINER ID        IMAGE                                       COMMAND                  CREATED             STATUS                     PORTS                              NAMES
dfcc0e58e111        cossacklabs/acra-server:latest              "/acra-server --conf…"   2 minutes ago       Up 2 minutes               9090/tcp, 0.0.0.0:9393->9393/tcp   acra-demo_acra-server_1
2601ddf7fb7b        postgres:11                                 "docker-entrypoint.s…"   2 minutes ago       Up 2 minutes               0.0.0.0:5432->5432/tcp             acra-demo_postgresql_1
9726f9355f56        cossacklabs/acra-keymaker:latest            "/acra-keymaker --cl…"   2 minutes ago       Exited (0) 2 minutes ago                                      acra-demo_acra-keymaker_server_1
ac0ca175f5be        cossacklabs/acra-poisonrecordmaker:latest   "./acra-poisonrecord…"   2 minutes ago       Exited (0) 2 minutes ago                                      acra-demo_acra-poisonrecordmaker_1
```

### Run demo app

Install dependencies and run demo application from repository folder. The demo application is [very simple](https://github.com/cossacklabs/acra-demo/blob/master/demo/demo.go), it works as database client application: connects to the database, creates test table, add some encrypted data, add poison records, reads data using `SELECT` query.

```bash
go run demo/demo.go -create
```

If no errors occurred, you should see log that table was created:

```sh
2022/03/19 13:53:44 Table has been successfully created
```

### Fill in the database table with data

Insert some data into the table by running client application again, for example, here we add 10 rows:

```bash
go run demo/demo.go -insert 10
```

If no errors, you should see:

```sh
2022/03/19 13:53:47 Insert has been successful
```

Client application adds some random data to the database but AcraServer sits transparently between app and database and encrypts all the data before storing in the database.

### Read the database table (aka steal all data)

Let’s check that we can read data from the table. Run client application with `-select` command.

```bash
go run demo/demo.go -select
```

If no errors, you should see all records:

```sh
1       yourheadisdewt8 P!RG3TVf+dL     fozkvhy@sylyznsg.abgtdgsx.org
2       peasmak1234k6   f2S38V5RO41     haqkhase@wggsqmqxuy.gfhcub.net
3       pagsulat08      A49vnzT9Sb*     oagbiufjxl@icjgsy.dcyygujs.org
4       madalianop      D962&6M7tsi     kyocxhgif@uqbdjx.tmsouz.com
5       shotgunroseswu  D962&6M7tsi     oqmzku@sqaxszyji.mzpyneh.com
6       anderswerdenod  giMOcrdOhz4     cqajxjvae@fsvqhhtu.omgvln.org
7       Buschborn6h     6G&8k-8R&_@     fozkvhy@sylyznsg.abgtdgsx.org
8       uvideli2d       D962&6M7tsi     vjnseokinl@aiiuujaoz.gwrfdfqm.net
9       macatsvy        1QH!k0_ZVk#     zchceeult@bvqpduqj.dnhdeg.com
10      anderswerdenod  LsRD%**0z2g     cgiefsi@qxshtw.tpweivkmzl.com
2022/03/19 14:04:03 Select has been successful
```

So we just downloaded all the content of the table. If we were attackers, we’d have successfully stolen all the data. As attackers, we could use some SQL injection to perform `SELECT *` query.

### Add poison records to prevent leak

Now we will add poison record to the table to detect an attack. Get the value of poison record data from the logs of exited `acra-poisonrecordmaker` container and then insert it into a table:

```bash
docker logs acra-demo-acra-poisonrecordmaker-1
```

If no errors, you should see base64 encoded value of poison record, it looks like encrypted data that we already have in the database (or like garbage):

```sh
IiIiIiIiIiJVRUMyAAAALaxV9EIC3i/fAgKysyzZUerLzfS17l72WKaFvnLidd8puf0xJfkgJwQmVAAAAAABAUAMAAAAEAAAACAAAABI15gRoCor8GWbMgamOioeaeZr149b/qk1LGpfSJ0+kHrtBNdP0rwKcdh0zsZgAnHZnRkXonklDDO4d4ZDAAAAAAAAAAABAUAMAAAAEAAAABcAAAAZDRadADVJbcS4CZ4hI0vAMh6em+Dy/B48xtoOfdWEQyYibGwDjUtp8pydV41ZQ91SU2U=
```

Copy poison record from your log of your `acra-poisonrecordmaker` and insert it to the database table:

```bash
go run demo/demo.go -insert_poison IiIiIiIiIiJVRUMyAAAALWSWDMcDH/+0AgCR2bsCZZW47bPtG+WtSD6Riq1PX/NxL1pCpeUgJwQmVAAAAAABAUAMAAAAEAAAACAAAABQeXSzlAcOIYtObhgHLTzGdCKFoEcoBJdtSjmxRtbTZplrFMQMTz15Ieww2FRBbSFN8sH0+pRmtjVxTEWEAAAAAAAAAAABAUAMAAAAEAAAAFgAAAB8UwNKO/MhI0ECetlJfELaqao/L1/WpvrEpGkol2h4MJIl4Mjo2CfEoAICOcJcbfeHPcKCCTtnUFgRhA4b0998U0j5bqBmmFvANHK0mPJMS37xWeLErxUtH/LgJ6ZdDYGg2/TkfS1+cxR/MLuJ93Nkrlf9VQ==
```

If no errors show, you should see base64 encoded value of poison record:

```sh
2022/03/19 14:05:42 Poison record insert has been successful
```

### Try to steal data again

Now we are protected from malicious `SELECT *` queries. Try to read all the data again:

```bash
go run demo/demo.go -select
```

You should see:

```sh
1       yourheadisdewt8 P!RG3TVf+dL     fozkvhy@sylyznsg.abgtdgsx.org
2       peasmak1234k6   f2S38V5RO41     haqkhase@wggsqmqxuy.gfhcub.net
3       pagsulat08      A49vnzT9Sb*     oagbiufjxl@icjgsy.dcyygujs.org
4       madalianop      D962&6M7tsi     kyocxhgif@uqbdjx.tmsouz.com
5       shotgunroseswu  D962&6M7tsi     oqmzku@sqaxszyji.mzpyneh.com
6       anderswerdenod  giMOcrdOhz4     cqajxjvae@fsvqhhtu.omgvln.org
7       Buschborn6h     6G&8k-8R&_@     fozkvhy@sylyznsg.abgtdgsx.org
8       uvideli2d       D962&6M7tsi     vjnseokinl@aiiuujaoz.gwrfdfqm.net
9       macatsvy        1QH!k0_ZVk#     zchceeult@bvqpduqj.dnhdeg.com
10      anderswerdenod  LsRD%**0z2g     cgiefsi@qxshtw.tpweivkmzl.com
2022/03/19 14:06:02 read tcp 127.0.0.1:51044->127.0.0.1:9393: read: connection reset by peer
exit status 1
```

Also, check the console where you run infrastructure. You should see that poison records has been detected by AcraServer:

```sh
acra-demo-acra-server-1             | time="2022-03-19T12:06:02Z" level=warning msg="Recognized poison record" client_id=acro_demo code=587 session_id=5
acra-demo-acra-server-1             | time="2022-03-19T12:06:02Z" level=warning msg="Recognized poison record"
acra-demo-acra-server-1             | time="2022-03-19T12:06:02Z" level=warning msg="Detected poison record, exit" code=101
acra-demo-acra-server-1 exited with code 1
```

It means that AcraServer detected poison record and stopped working (shut down itself). Note, that you can setup AcraServer to different behaviour when it detects poison record:

- perform a shut-down (useful for very critical data, but AcraServer will be down until you restart it)
- run a script (you can tell AcraServer to run a script after detecting poison record, for example, to send alerts to system administrators and SIEMs)
- perform a shut-down and run a script

## Further steps

Let us know if you have any questions by dropping an email to [dev@cossacklabs.com](mailto:dev@cossacklabs.com).

1. [Acra features](https://cossacklabs.com/acra/) – check out full features set and available licenses.
2. Other [Acra example applications](https://github.com/cossacklabs/acra-engineering-demo/) – try other Acra features, like transparent encryption, SQL firewall, load balancing support.

## Need help?

Need help in configuring Acra? Our support is available for [Acra Pro and Acra Enterprise versions](https://www.cossacklabs.com/acra/#pricing).
