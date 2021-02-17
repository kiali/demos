package demo.runtimes.springboot;

import com.sun.tools.javac.util.List;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.reactive.function.client.WebClient;

@Controller
public class Greetings {

	@GetMapping("/greetings")
	@ResponseBody
	public String greetings(@RequestParam(name="from", required=false) String from) {
		if (from == null) {
			List.of("vertx", "quarkus", "nodejs").forEach(peer -> {
				WebClient.create("http://" + peer + ":8080")
					.get()
					.uri("/greetings?from=springboot")
					.retrieve()
					.bodyToMono(String.class)
					.doOnNext(System.out::println);
			});
			String resp = "Hello, I'm SpringBoot!";
			System.out.println(resp);
			return resp;
		} else {
			String resp = "Hello " + from + ", I'm SpringBoot!";
			System.out.println(resp);
			return resp;
		}
	}
}
