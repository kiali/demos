package demo.runtimes.vertx

import io.micrometer.core.instrument.binder.jvm.ClassLoaderMetrics
import io.micrometer.core.instrument.binder.jvm.JvmMemoryMetrics
import io.micrometer.core.instrument.binder.jvm.JvmThreadMetrics
import io.vertx.core.Vertx
import io.vertx.core.VertxOptions
import io.vertx.ext.web.client.WebClient
import io.vertx.micrometer.Label
import io.vertx.micrometer.MicrometerMetricsOptions
import io.vertx.micrometer.VertxPrometheusOptions
import io.vertx.micrometer.backends.BackendRegistries
import java.util.*

/**
 * @author Joel Takvorian
 */
fun main() {
  val vertx = vertx()
  val client = WebClient.create(vertx)
  vertx.deployVerticle(Greetings(client))
}

fun vertx(): Vertx {
  val opts = VertxOptions().setMetricsOptions(
          MicrometerMetricsOptions()
                  .setPrometheusOptions(VertxPrometheusOptions()
                          .setPublishQuantiles(true)
                          .setEnabled(true))
                  .setLabels(EnumSet.of(Label.POOL_TYPE, Label.POOL_NAME, Label.CLASS_NAME, Label.HTTP_CODE, Label.HTTP_METHOD, Label.HTTP_PATH, Label.EB_ADDRESS, Label.EB_FAILURE, Label.EB_SIDE))
                  .setEnabled(true)
  )
  val vertx = Vertx.vertx(opts)
  // Instrument JVM
  val registry = BackendRegistries.getDefaultNow()
  if (registry != null) {
    ClassLoaderMetrics().bindTo(registry)
    JvmMemoryMetrics().bindTo(registry)
    // JvmGcMetrics().bindTo(registry);
    // ProcessorMetrics().bindTo(registry);
    JvmThreadMetrics().bindTo(registry)
  }
  return vertx.exceptionHandler { it.printStackTrace() }
}