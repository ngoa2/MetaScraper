docker rm -f ngoa2server

docker pull ngoa2/ngoa2server

docker run -d \
-p 443:443 \
--name ngoa2server \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSKEY=/etc/letsencrypt/live/api.ngoa2.me/privkey.pem \
-e TLSCERT=/etc/letsencrypt/live/api.ngoa2.me/fullchain.pem ngoa2/ngoa2server

exit