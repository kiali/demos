FROM ubi7/nodejs-12

EXPOSE 8080

COPY nodejs /app/

CMD cd /app && npm start
