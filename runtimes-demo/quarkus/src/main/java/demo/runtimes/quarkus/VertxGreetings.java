package demo.runtimes.quarkus;

import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.QueryParam;

@Path("/")
@RegisterRestClient
public interface VertxGreetings {
  @GET
  @Path("/greetings")
  String greetings(@QueryParam("from") String from);
}
