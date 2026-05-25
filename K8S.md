# Kubernetes Production Readiness

## Scaling

Das Deployment `product-catalog-api` wurde mit folgendem Befehl auf drei Replikas skaliert:

```bash
kubectl scale deployment product-catalog-api --replicas=3 -n product-catalog
```

Anschließend wurde mit kubectl get pods -n product-catalog überprüft, dass drei API-Pods erstellt wurden.

## Readiness vs. Liveness Probe
Eine Readiness Probe prüft, ob ein Pod bereit ist, Traffic entgegenzunehmen. Wenn die Readiness Probe fehlschlägt, bleibt der Pod zwar gestartet, wird aber vorübergehend aus dem Service entfernt und erhält keine Requests.
Eine Liveness Probe prüft, ob die Anwendung im Container noch lebt. Wenn die Liveness Probe fehlschlägt, wird der Container von Kubernetes neu gestartet.
Unterschiedliche initialDelaySeconds Werte sind sinnvoll, weil eine Anwendung oft zuerst Zeit zum Starten benötigt, bevor sie Requests zuverlässig beantworten kann. Die Liveness Probe sollte meist etwas später starten, damit Kubernetes den Container nicht zu früh neu startet.

## Resource Requests and Limits
Resource Requests geben an, wie viel CPU und Memory ein Container mindestens benötigt. Kubernetes nutzt diese Werte für das Scheduling, also um einen passenden Node auszuwählen.
Resource Limits legen fest, wie viele Ressourcen ein Container maximal verwenden darf. Wenn ein Container das Memory Limit überschreitet, kann er beendet werden. Wenn er das CPU Limit überschreitet, wird seine CPU-Nutzung gedrosselt.
Requests und Limits sollten beide gesetzt werden, damit Kubernetes Ressourcen besser planen kann und einzelne Pods nicht zu viele Ressourcen verbrauchen.
