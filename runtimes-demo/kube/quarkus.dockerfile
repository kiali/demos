FROM fabric8/java-alpine-openjdk11-jre
EXPOSE 8080
ENV JAVA_OPTIONS="-Dquarkus.http.host=0.0.0.0 -Djava.util.logging.manager=org.jboss.logmanager.LogManager"
ENV AB_ENABLED=jmx_exporter
COPY quarkus/target/lib/* /deployments/lib/
COPY quarkus/target/*-runner.jar /deployments/app.jar
ENTRYPOINT [ "/deployments/run-java.sh" ]
