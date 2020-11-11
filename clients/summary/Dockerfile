# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base

FROM NGINX

ADD default.conf /etc/nginx/conf.d/default.conf
COPY client client
COPY index index

