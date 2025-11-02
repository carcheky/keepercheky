/**
 * Dashboard Charts - Chart.js integration for KeeperCheky Dashboard
 * 
 * This file contains all chart creation and management functions for the enhanced dashboard.
 * Uses Chart.js 4.4.0 with dark theme optimizations.
 */

// Chart instances registry for cleanup
window.dashboardCharts = {
    diskUsage: null,
    mediaDistribution: null,
    qualityDistribution: null
};

/**
 * Creates a line chart showing disk usage evolution over 30 days
 * @param {string} canvasId - Canvas element ID
 * @param {Array} data - Array of DiskUsagePoint objects
 */
function createDiskUsageChart(canvasId, data) {
    // Destroy existing chart if any
    if (window.dashboardCharts.diskUsage) {
        window.dashboardCharts.diskUsage.destroy();
    }

    const ctx = document.getElementById(canvasId);
    if (!ctx) {
        console.warn('Canvas element not found:', canvasId);
        return;
    }

    window.dashboardCharts.diskUsage = new Chart(ctx.getContext('2d'), {
        type: 'line',
        data: {
            labels: data.map(point => {
                const date = new Date(point.date);
                return date.toLocaleDateString('es-ES', { month: 'short', day: 'numeric' });
            }),
            datasets: [
                {
                    label: 'Espacio Usado (GB)',
                    data: data.map(point => point.used_gb.toFixed(2)),
                    borderColor: 'rgb(59, 130, 246)',
                    backgroundColor: 'rgba(59, 130, 246, 0.1)',
                    tension: 0.4,
                    fill: true,
                    pointRadius: 2,
                    pointHoverRadius: 5
                },
                {
                    label: 'Espacio Libre (GB)',
                    data: data.map(point => point.free_gb.toFixed(2)),
                    borderColor: 'rgb(34, 197, 94)',
                    backgroundColor: 'rgba(34, 197, 94, 0.1)',
                    tension: 0.4,
                    fill: true,
                    pointRadius: 2,
                    pointHoverRadius: 5
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
                intersect: false,
                mode: 'index'
            },
            plugins: {
                title: {
                    display: true,
                    text: 'Evolución del Espacio en Disco (30 días)',
                    color: '#e2e8f0',
                    font: {
                        size: 16,
                        weight: 'bold'
                    },
                    padding: {
                        top: 10,
                        bottom: 20
                    }
                },
                legend: {
                    labels: { 
                        color: '#e2e8f0',
                        padding: 15,
                        font: {
                            size: 12
                        }
                    }
                },
                tooltip: {
                    backgroundColor: 'rgba(30, 41, 59, 0.95)',
                    titleColor: '#e2e8f0',
                    bodyColor: '#e2e8f0',
                    borderColor: '#334155',
                    borderWidth: 1,
                    padding: 12,
                    displayColors: true,
                    callbacks: {
                        label: function(context) {
                            return context.dataset.label + ': ' + parseFloat(context.parsed.y).toFixed(2) + ' GB';
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: { 
                        color: '#94a3b8',
                        callback: function(value) {
                            return value.toFixed(0) + ' GB';
                        }
                    },
                    grid: { 
                        color: 'rgba(255, 255, 255, 0.05)'
                    }
                },
                x: {
                    ticks: { 
                        color: '#94a3b8',
                        maxRotation: 45,
                        minRotation: 0
                    },
                    grid: { 
                        color: 'rgba(255, 255, 255, 0.05)'
                    }
                }
            }
        }
    });
}

/**
 * Creates a doughnut chart showing media distribution by type
 * @param {string} canvasId - Canvas element ID
 * @param {Object} stats - Stats object with total_movies and total_series
 */
function createMediaDistributionChart(canvasId, stats) {
    // Destroy existing chart if any
    if (window.dashboardCharts.mediaDistribution) {
        window.dashboardCharts.mediaDistribution.destroy();
    }

    const ctx = document.getElementById(canvasId);
    if (!ctx) {
        console.warn('Canvas element not found:', canvasId);
        return;
    }

    const movies = stats.total_movies || 0;
    const series = stats.total_series || 0;
    const total = movies + series;

    window.dashboardCharts.mediaDistribution = new Chart(ctx.getContext('2d'), {
        type: 'doughnut',
        data: {
            labels: ['Películas', 'Series'],
            datasets: [{
                data: [movies, series],
                backgroundColor: [
                    'rgba(59, 130, 246, 0.8)',   // Blue
                    'rgba(168, 85, 247, 0.8)'    // Purple
                ],
                borderColor: [
                    'rgba(59, 130, 246, 1)',
                    'rgba(168, 85, 247, 1)'
                ],
                borderWidth: 2,
                hoverOffset: 10
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Distribución por Tipo',
                    color: '#e2e8f0',
                    font: {
                        size: 16,
                        weight: 'bold'
                    },
                    padding: {
                        top: 10,
                        bottom: 20
                    }
                },
                legend: {
                    position: 'bottom',
                    labels: { 
                        color: '#e2e8f0',
                        padding: 15,
                        font: {
                            size: 13
                        },
                        generateLabels: function(chart) {
                            const data = chart.data;
                            if (data.labels.length && data.datasets.length) {
                                return data.labels.map((label, i) => {
                                    const value = data.datasets[0].data[i];
                                    const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : 0;
                                    return {
                                        text: `${label}: ${value} (${percentage}%)`,
                                        fillStyle: data.datasets[0].backgroundColor[i],
                                        hidden: false,
                                        index: i
                                    };
                                });
                            }
                            return [];
                        }
                    }
                },
                tooltip: {
                    backgroundColor: 'rgba(30, 41, 59, 0.95)',
                    titleColor: '#e2e8f0',
                    bodyColor: '#e2e8f0',
                    borderColor: '#334155',
                    borderWidth: 1,
                    padding: 12,
                    callbacks: {
                        label: function(context) {
                            const value = context.parsed;
                            const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : 0;
                            return context.label + ': ' + value + ' (' + percentage + '%)';
                        }
                    }
                }
            }
        }
    });
}

/**
 * Creates a horizontal bar chart showing media distribution by quality
 * @param {string} canvasId - Canvas element ID
 * @param {Object} distribution - Object mapping quality to count
 */
function createQualityDistributionChart(canvasId, distribution) {
    // Destroy existing chart if any
    if (window.dashboardCharts.qualityDistribution) {
        window.dashboardCharts.qualityDistribution.destroy();
    }

    const ctx = document.getElementById(canvasId);
    if (!ctx) {
        console.warn('Canvas element not found:', canvasId);
        return;
    }

    // Sort qualities by count descending
    const qualities = Object.keys(distribution).sort((a, b) => distribution[b] - distribution[a]);
    const counts = qualities.map(q => distribution[q]);

    // Generate colors based on quality
    const backgroundColors = qualities.map(quality => {
        if (quality.includes('2160p') || quality.includes('4K')) {
            return 'rgba(251, 146, 60, 0.8)'; // Orange for 4K
        } else if (quality.includes('1080p')) {
            return 'rgba(34, 197, 94, 0.8)'; // Green for 1080p
        } else if (quality.includes('720p')) {
            return 'rgba(59, 130, 246, 0.8)'; // Blue for 720p
        } else {
            return 'rgba(148, 163, 184, 0.8)'; // Gray for others
        }
    });

    // Border colors with full opacity
    const borderColors = qualities.map(quality => {
        if (quality.includes('2160p') || quality.includes('4K')) {
            return 'rgba(251, 146, 60, 1)'; // Orange for 4K
        } else if (quality.includes('1080p')) {
            return 'rgba(34, 197, 94, 1)'; // Green for 1080p
        } else if (quality.includes('720p')) {
            return 'rgba(59, 130, 246, 1)'; // Blue for 720p
        } else {
            return 'rgba(148, 163, 184, 1)'; // Gray for others
        }
    });

    window.dashboardCharts.qualityDistribution = new Chart(ctx.getContext('2d'), {
        type: 'bar',
        data: {
            labels: qualities,
            datasets: [{
                label: 'Archivos',
                data: counts,
                backgroundColor: backgroundColors,
                borderColor: borderColors,
                borderWidth: 1,
                borderRadius: 4
            }]
        },
        options: {
            indexAxis: 'y', // Horizontal bars
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Distribución por Calidad',
                    color: '#e2e8f0',
                    font: {
                        size: 16,
                        weight: 'bold'
                    },
                    padding: {
                        top: 10,
                        bottom: 20
                    }
                },
                legend: { 
                    display: false 
                },
                tooltip: {
                    backgroundColor: 'rgba(30, 41, 59, 0.95)',
                    titleColor: '#e2e8f0',
                    bodyColor: '#e2e8f0',
                    borderColor: '#334155',
                    borderWidth: 1,
                    padding: 12,
                    callbacks: {
                        label: function(context) {
                            return 'Archivos: ' + context.parsed.x;
                        }
                    }
                }
            },
            scales: {
                x: {
                    beginAtZero: true,
                    ticks: { 
                        color: '#94a3b8',
                        precision: 0
                    },
                    grid: { 
                        color: 'rgba(255, 255, 255, 0.05)'
                    }
                },
                y: {
                    ticks: { 
                        color: '#94a3b8',
                        font: {
                            size: 11
                        }
                    },
                    grid: { 
                        display: false 
                    }
                }
            }
        }
    });
}

/**
 * Formats timestamp to human-readable relative time in Spanish
 * @param {string} timestampStr - ISO timestamp string
 * @returns {string} Formatted time string
 */
window.formatTimestamp = function formatTimestamp(timestampStr) {
    if (!timestampStr) return '';
    
    const timestamp = new Date(timestampStr);
    const now = new Date();
    const diffMs = now - timestamp;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);
    
    if (diffMins < 1) return 'Ahora mismo';
    if (diffMins < 60) return `Hace ${diffMins} min`;
    if (diffHours < 24) return `Hace ${diffHours}h`;
    if (diffDays === 1) return 'Ayer';
    if (diffDays < 7) return `Hace ${diffDays} días`;
    
    return timestamp.toLocaleDateString('es-ES', { 
        month: 'short', 
        day: 'numeric',
        year: timestamp.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
    });
};

/**
 * Cleanup all dashboard charts when component is destroyed
 */
function destroyDashboardCharts() {
    Object.values(window.dashboardCharts).forEach(chart => {
        if (chart) {
            chart.destroy();
        }
    });
    window.dashboardCharts = {
        diskUsage: null,
        mediaDistribution: null,
        qualityDistribution: null
    };
}
