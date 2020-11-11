docker rm -f ngoa2client

docker pull ngoa2/ngoa2client

docker run -d \
-p 443:443 \
-p 80:80 \
--name ngoa2client \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSKEY=/etc/letsencrypt/live/ngoa2.me/privkey.pem \
-e TLSCERT=/etc/letsencrypt/live/ngoa2.me/fullchain.pem ngoa2/ngoa2client

exit