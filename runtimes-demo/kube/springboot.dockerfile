FROM adoptopenjdk/openjdk11:alpine-slim

EXPOSE 8080

# Copy dependencies
COPY springboot/target/dependency/* /deployment/libs/

# Copy classes
COPY springboot/target/classes /deployment/classes

RUN chgrp -R 0 /deployment && chmod -R g+rwX /deployment

CMD java -cp /deployment/classes:/deployment/libs/* demo.runtimes.springboot.GreetingsApplication
