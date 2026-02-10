package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusController struct {
	healthController *HealthController
}

func NewStatusController() *StatusController {
	return &StatusController{
		healthController: NewHealthController(),
	}
}

// @Summary Status Dashboard
// @Description Displays a real-time HTML status dashboard for monitoring system health
// @Tags Health
// @Produce html
// @Success 200 {string} string "HTML status dashboard"
// @Router /status [get]
func (c *StatusController) ShowDashboard(ctx *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Status</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #fafafa;
            color: #333;
        }

        .container {
            max-width: 900px;
            margin: 0 auto;
            padding: 40px 20px;
        }

        .header {
            margin-bottom: 40px;
        }

        .header h1 {
            font-size: 2rem;
            font-weight: 600;
            margin-bottom: 8px;
        }

        .header .time {
            color: #666;
            font-size: 0.9rem;
        }

        .status-banner {
            background: white;
            border-left: 4px solid #22c55e;
            padding: 20px;
            margin-bottom: 30px;
            border-radius: 4px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }

        .status-banner.degraded {
            border-left-color: #f59e0b;
        }

        .status-banner.unhealthy {
            border-left-color: #ef4444;
        }

        .status-text {
            font-size: 1.1rem;
            font-weight: 500;
        }

        .status-meta {
            color: #666;
            font-size: 0.85rem;
            margin-top: 8px;
        }

        .services {
            background: white;
            border-radius: 4px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .service-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 16px 20px;
            border-bottom: 1px solid #f0f0f0;
        }

        .service-row:last-child {
            border-bottom: none;
        }

        .service-name {
            font-weight: 500;
        }

        .service-info {
            color: #666;
            font-size: 0.85rem;
            margin-top: 4px;
        }

        .service-status {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 0.9rem;
        }

        .status-dot {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background: #22c55e;
        }

        .status-dot.down {
            background: #ef4444;
        }

        .refresh-notice {
            text-align: center;
            color: #999;
            font-size: 0.85rem;
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>System Status</h1>
            <div class="time" id="timestamp">Loading...</div>
        </div>

        <div class="status-banner" id="status-banner">
            <div class="status-text" id="status-text">Loading...</div>
            <div class="status-meta" id="status-meta"></div>
        </div>

        <div class="services" id="services">
            <div class="service-row">
                <div>Loading services...</div>
            </div>
        </div>

        <div class="refresh-notice">
            Auto-refreshes every 5 seconds
        </div>
    </div>

    <script>
        async function fetchStatus() {
            try {
                const res = await fetch('/health/detailed');
                const data = await res.json();
                updateUI(data);
            } catch (error) {
                document.getElementById('status-text').textContent = 'Failed to load status';
            }
        }

        function updateUI(data) {
            const time = new Date(data.timestamp).toLocaleString('en-US', {
                month: 'short',
                day: 'numeric',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit'
            });
            document.getElementById('timestamp').textContent = time;

            const banner = document.getElementById('status-banner');
            banner.className = 'status-banner ' + data.status;

            const statusMessages = {
                'healthy': 'All systems operational',
                'degraded': 'Partial system outage',
                'unhealthy': 'System outage'
            };

            document.getElementById('status-text').textContent = statusMessages[data.status] || data.status;
            document.getElementById('status-meta').textContent = 'Uptime: ' + data.uptime + ' · Version ' + data.version;

            let html = '';
            for (const [name, service] of Object.entries(data.services)) {
                const statusClass = service.status === 'up' ? '' : 'down';
                const statusText = service.status === 'up' ? 'Operational' : 'Down';

                let info = service.message || '';
                if (service.latency_ms) {
                    info += (info ? ' · ' : '') + service.latency_ms.toFixed(1) + 'ms';
                }
                if (name === 'database' && service.details) {
                    info += (info ? ' · ' : '') + 'Connections: ' + service.details.active + '/' + service.details.max;
                }

                html += '<div class="service-row">' +
                    '<div>' +
                        '<div class="service-name">' + name.charAt(0).toUpperCase() + name.slice(1) + '</div>' +
                        (info ? '<div class="service-info">' + info + '</div>' : '') +
                    '</div>' +
                    '<div class="service-status">' +
                        '<span class="status-dot ' + statusClass + '"></span>' +
                        '<span>' + statusText + '</span>' +
                    '</div>' +
                '</div>';
            }

            document.getElementById('services').innerHTML = html;
        }

        fetchStatus();
        setInterval(fetchStatus, 5000);
    </script>
</body>
</html>`

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
