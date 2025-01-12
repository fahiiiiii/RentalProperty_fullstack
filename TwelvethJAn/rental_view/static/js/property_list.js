function renderProperties(properties) {
    const container = document.getElementById('listingsContainer');
    container.innerHTML = properties.map(property => `
        <div class="listing-card" data-property-id="${property.id}">
            <div class="property-image-container">
                <img 
                    src="/static/images/${property.images[0] || 'placeholder.jpg'}" 
                    alt="${property.title}" 
                    class="property-image"
                >
                <button class="favorite-btn" onclick="toggleFavorite(${property.id})">
                    <i class="heart-icon">♡</i>
                </button>
            </div>
            <div class="property-details">
                <div class="property-header">
                    <h3 class="property-title">${property.title}</h3>
                    <div class="property-rating">
                        <span class="rating-value">${property.rating.toFixed(1)}</span>
                        <span class="rating-reviews">(${property.reviews} reviews)</span>
                    </div>
                </div>
                <div class="property-meta">
                    <div class="property-location">
                        <i class="location-icon">📍</i>
                        ${property.location}
                    </div>
                    <div class="property-type">
                        <i class="type-icon">🏠</i>
                        ${property.type}
                    </div>
                </div>
                <div class="property-features">
                    <div class="feature">
                        <i class="bed-icon">🛏️</i>
                        ${property.bedrooms} Bedrooms
                    </div>
                    <div class="feature">
                        <i class="bath-icon">🚿</i>
                        ${property.bathrooms} Bathrooms
                    </div>
                </div>
                <div class="property-amenities">
                    ${property.amenities.slice(0, 3).map(amenity => `
                        <span class="amenity-tag">${amenity}</span>
                    `).join('')}
                    ${property.amenities.length > 3 ? 
                        `<span class="amenity-more">+${property.amenities.length - 3} more</span>` : 
                        ''
                    }
                </div>
                <div class="property-pricing">
                    <div class="price">
                        $${property.price.toLocaleString()} 
                        <span class="price-period">/night</span>
                    </div>
                    <button 
                        class="book-btn" 
                        onclick="viewPropertyDetails(${property.id})"
                    >
                        View Details
                    </button>
                </div>
            </div>
        </div>
    `).join('');
}

function updatePagination(pagination) {
    const paginationContainer = document.getElementById('paginationContainer');
    const { page, totalPages } = pagination;

    // Clear existing pagination
    paginationContainer.innerHTML = '';

    // Previous page button
    if (page > 1) {
        const prevBtn = document.createElement('button');
        prevBtn.textContent = 'Previous';
        prevBtn.onclick = () => fetchProperties(page - 1);
        paginationContainer.appendChild(prevBtn);
    }

    // Page numbers
    for (let i = 1; i <= totalPages; i++) {
        const pageBtn = document.createElement('button');
        pageBtn.textContent = i;
        pageBtn.classList.toggle('active', i === page);
        pageBtn.onclick = () => fetchProperties(i);
        paginationContainer.appendChild(pageBtn);
    }

    // Next page button
    if (page < totalPages) {
        const nextBtn = document.createElement('button');
        nextBtn.textContent = 'Next';
        nextBtn.onclick = () => fetchProperties(page + 1);
        paginationContainer.appendChild(nextBtn);
    }
}

function toggleFavorite(propertyId) {
    const favoriteBtn = document.querySelector(
        `.listing-card[data-property-id="${propertyId}"] .favorite-btn`
    );
    
    favoriteBtn.classList.toggle('favorited');
    
    // Optional: Send favorite status to backend
    // sendFavoriteStatus(propertyId, favoriteBtn.classList.contains('favorited'));
}

function viewPropertyDetails(propertyId) {
    // Redirect or open modal with property details
    window.location.href = `/property/${propertyId}`;
}

// Advanced filtering
function applyFilters() {
    const filters = {
        city: document.getElementById('cityFilter').value,
        minPrice: document.getElementById('minPriceFilter').value,
        maxPrice: document.getElementById('maxPriceFilter').value,
        bedrooms: document.getElementById('bedroomsFilter').value,
        amenities: Array.from(
            document.querySelectorAll('input[name="amenities"]:checked')
        ).map(checkbox => checkbox.value)
    };

    fetchFilteredProperties(filters);
}

async function fetchFilteredProperties(filters) {
    try {
        const response = await fetch('/v1/property/filter', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(filters)
        });

        const result = await response.json();

        if (result.success) {
            renderProperties(result.data);
        } else {
            console.error('Filtering failed:', result.error);
        }
    } catch (error) {
        console.error('Error applying filters:', error);
    }
}

// Initial load
document.addEventListener('DOMContentLoaded', () => {
    fetchProperties();
});