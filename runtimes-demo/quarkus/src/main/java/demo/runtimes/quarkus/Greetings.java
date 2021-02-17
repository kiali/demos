package demo.runtimes.quarkus;

import org.eclipse.microprofile.rest.client.inject.RestClient;

import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;

@Path("/")
public class Greetings {

    // Ah, frameworks ¯\_(ツ)_/¯

    @Inject
    @RestClient
    VertxGreetings vertx;

    @Inject
    @RestClient
    NodejsGreetings nodejs;

    @Inject
    @RestClient
    GoGreetings go;

    @Inject
    @RestClient
    SpringbootGreetings springboot;

    @GET
    @Path("/greetings")
    @Produces("text/plain")
    public String greetings(@QueryParam("from") String from) {
        if (from == null) {
            String resp = "Hello, I'm Quarkus!";
            String me = "quarkus";
            try {
                resp += "\n" + vertx.greetings(me);
            } catch (Exception e) {
                e.printStackTrace();
            }
            try {
                resp += "\n" + nodejs.greetings(me);
            } catch (Exception e) {
                e.printStackTrace();
            }
            try {
                resp += "\n" + go.greetings(me);
            } catch (Exception e) {
                e.printStackTrace();
            }
            try {
                resp += "\n" + springboot.greetings(me);
            } catch (Exception e) {
                e.printStackTrace();
            }
            System.out.println(resp);
            return resp;
        } else {
            String resp = "Hello " + from + ", I'm Quarkus!";
            System.out.println(resp);
            return resp;
        }
    }
}
