
import java.time.Duration;
import java.time.Instant;
import java.util.List;
import javax.ws.rs.client.Client;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.client.WebTarget;
import javax.ws.rs.core.GenericType;

public class Main {
    public static void main(String[] args) {
        System.out.println("success/s\tfail/s");
        int success = 0, failure = 0;
        Instant start = Instant.now(), logged = Instant.now();
        Client client = ClientBuilder.newClient();
        while (Duration.between(start, Instant.now()).toMillis() < 10_000) {
            WebTarget target = client.target("https://george-json-test.s3.amazonaws.com/example.json");
            List<FooBar> resp = target.request().get(new GenericType<List<FooBar>>() {});
            if (ok(resp)) {
                success++;
            } else {
                failure++;
            }
            if (Duration.between(logged, Instant.now()).toMillis() > 1_000) {
                Duration runTime = Duration.between(logged, Instant.now());
                double successRate = success * 1000.0 / runTime.toMillis();
                double failureRate = failure * 1000.0 / runTime.toMillis();
                System.out.println(String.format("%d\t%d", (int) successRate, (int) failureRate));
                success = 0;
                failure = 0;
                logged = Instant.now();
            }
        }
    }

    private static boolean ok(List<FooBar> resp) {
        for (FooBar row : resp) {
            if (row.foo != 1 || row.bar != 2) {
                return false;
            }
        }
        return true;
    }

    static class FooBar {
        public int foo, bar;
    }
}