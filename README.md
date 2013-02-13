nettis
=======

nettis is an integration testing utility for network apps.

nettis echos back data as a TCP server or client. Optionally use HTTP(S) &amp; TLS

NOTE: nettis is only a testing tool. Only use it from behind a firewall.

Synopsis
-------
nettis [options] [host:]<port>

### Options

>  -d=0: delay (seconds) before echoing
>  -h=false: Show this help
>  -http=false: use http
>  -i=false: initiate conversation
>  -l=false: listen
>  -s=false: Secure sockets (TLS/SSL)
>  -s-cert="cert.pem": Certificate to use for TLS
>  -s-key="key.pem": Key to use for TLS
>  -s-trusted-cert="": Trusted certificate to accept TLS (nil means trust-all)
>  -v=false: verbose
>  -version=false: Show version

Examples
--------
Listen on a port and echo back anything which comes in:

      nettis -l 9000

Connect on to the same port. Initiate a conversation with the server:

      nettis -i 9000

Start an HTTPS server (which echoes back request bodies):

      nettis -http -s 9443

Note: nettis automatically generates certificates & keys if they're not already there.

To do
-----

 * more granular verbosity
 * option to specify 'message content', particularly in conjunction with -i
