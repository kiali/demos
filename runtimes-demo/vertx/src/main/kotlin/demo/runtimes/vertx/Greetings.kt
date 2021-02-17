package demo.runtimes.vertx

import io.vertx.core.http.HttpServerOptions
import io.vertx.ext.web.Router
import io.vertx.ext.web.RoutingContext
import io.vertx.ext.web.client.WebClient
import io.vertx.kotlin.coroutines.CoroutineVerticle
import io.vertx.micrometer.PrometheusScrapingHandler

class Greetings(private val client: WebClient) : CoroutineVerticle() {
  private val peers = listOf("quarkus", "nodejs", "springboot")

  override suspend fun start() {
    // Register stadium API
    val serverOptions = HttpServerOptions().setPort(8080)
    val router = Router.router(vertx)
    router.route("/metrics").handler(PrometheusScrapingHandler.create())

    router["/health"].handler { it.response().end() }
    router["/greetings"].handler { greetings(it) }

    vertx.createHttpServer().requestHandler(router)
            .listen(serverOptions.port, serverOptions.host)
  }

  private fun greetings(ctx: RoutingContext) {
    val from = ctx.request().getParam("from")
    if (from == null) {
      ctx.response().end("Hello, I'm Vert.X!")
      peers.forEach { peer ->
        client.get(8080, peer, "/greetings?from=vertx").send {
          if (it.succeeded()) {
            println(it.result().bodyAsString())
          } else {
            it.cause().printStackTrace()
          }
        }
      }
    } else {
      val resp = "Hello $from, I'm Vert.X!"
      println(resp)
      ctx.response().end(resp)
    }
  }
}
