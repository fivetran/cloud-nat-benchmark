
import java.time.Duration;
import java.time.Instant;
import javax.ws.rs.client.Client;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.client.WebTarget;

public class Main {
    public static void main(String[] args) {
        System.out.println("success/s\tfail/s\t");
        int success = 0, failure = 0;
        Instant start = Instant.now(), logged = Instant.now();
        Client client = ClientBuilder.newClient();
        while (Duration.between(start, Instant.now()).toMillis() < 10_000) {
            WebTarget target = client.target("http://example.com");
            String resp = target.request().get(String.class);
            if (resp.hashCode() != -801093019) {
                System.err.println(resp.hashCode());
                failure++;
            } else {
                success++;
            }
            if (Duration.between(logged, Instant.now()).toMillis() > 1_000) {
                Duration runTime = Duration.between(logged, Instant.now());
                double successRate = success * 1000.0 / runTime.toMillis();
                double failureRate = failure * 1000.0 / runTime.toMillis();
                System.out.println(String.format("%d\t%d\n", (int) successRate, (int) failureRate));
                success = 0;
                failure = 0;
                logged = Instant.now();
            }
        }
    }
}