FROM adoptopenjdk/openjdk11:alpine-slim

EXPOSE 8080

# Copy dependencies
COPY vertx/target/dependency/* /deployment/libs/

# Copy classes
COPY vertx/target/classes /deployment/classes

RUN chgrp -R 0 /deployment && chmod -R g+rwX /deployment

CMD java -Dvertx.disableDnsResolver=true -cp /deployment/classes:/deployment/libs/* demo.runtimes.vertx.MainKt
