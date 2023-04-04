FROM nginx:latest

RUN apt-get update && \
    apt-get install -y clamav clamav-freshclam cron

RUN mkdir -p /usr/share/nginx/html/clamav

# Change ownership of the clamav directory to the clamav user and group
RUN chown clamav:clamav /usr/share/nginx/html/clamav

# Download the initial signature files so they're present once already when the container is up
RUN freshclam --datadir=/usr/share/nginx/html/clamav

# Update the Nginx configuration to serve the signature files
RUN echo 'server {' > /etc/nginx/conf.d/clamav.conf && \
    echo '    listen 80;' >> /etc/nginx/conf.d/clamav.conf && \
    echo '    location /clamav/ {' >> /etc/nginx/conf.d/clamav.conf && \
    echo '        alias /usr/share/nginx/html/clamav/;' >> /etc/nginx/conf.d/clamav.conf && \
    echo '        autoindex on;' >> /etc/nginx/conf.d/clamav.conf && \
    echo '    }' >> /etc/nginx/conf.d/clamav.conf && \
    echo '}' >> /etc/nginx/conf.d/clamav.conf

# Set up a cron job to update the database files every 3 hours
RUN echo "0 */3 * * * root freshclam --datadir=/usr/share/nginx/html/clamav --quiet" > /etc/cron.d/freshclam-mirror

# Make sure the cron job file has proper permissions
RUN chmod 0644 /etc/cron.d/freshclam-mirror

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Use the entrypoint script to start freshclam, cron, and Nginx
ENTRYPOINT ["/entrypoint.sh"]
