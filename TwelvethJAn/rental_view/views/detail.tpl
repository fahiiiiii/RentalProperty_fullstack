<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.listing.Title}}</title>
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body>
    <div class="container">
        <div class="property-header">
            <h1 class="property-title">{{.listing.Title}}</h1>
            <div class="property-stats">
                <span class="rating">★★★★★ {{.listing.Rating}}</span>
                <span>({{.listing.Reviews}} Reviews)</span>
                <div class="stats-item">
                    <span>🛏 {{.listing.Bedrooms}} Bedroom</span>
                </div>
                <div class="stats-item">
                    <span>🚽 {{.listing.Bathrooms}} Bathroom</span>
                </div>
                <div class="stats-item">
                    <span>👥 {{.listing.Guests}} Guests</span>
                </div>
            </div>
        </div>
        
        <div class="gallery">
            <img src="/static/images/{{index .listing.Images 0}}" alt="{{.listing.Title}}" class="main-image">
            <div class="thumbnail-grid">
                {{range .listing.Images}}
                    <img src="/static/images/{{.}}" alt="{{$.listing.Title}}" class="thumbnail">
                {{end}}
            </div>
        </div>
        
        <div class="property-info">
            <h2>{{.listing.Title}}</h2>
            <p class="property-description">
                {{.listing.Description}}
            </p>
        </div>
    </div>
</body>
</html>