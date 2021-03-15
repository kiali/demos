FROM centos/nodejs-12-centos7

EXPOSE 8080

COPY nodejs /app/

CMD cd /app && npm start
