// Sample data structure for the hostel
const hostelInfo = {
    title: "Hostel - 8 guests | Hostel in Dubai with Pool",
    rating: "★★★★★ 10.0",
    reviews: 520,
    stats: {
        bedrooms: 1,
        bathrooms: 1,
        guests: 8,
    },
    description: [
        "Welcome to Papaya's Backpackers, the hidden gem of Dubai...",
        // Additional descriptions
    ],
    images: {
        main: "pool-view.jpg",
        thumbnails: [
            "common-area-1.jpg",
            "hallway.jpg",
            "kitchen.jpg",
            "lounge.jpg",
        ],
    },
};

exports.getHostelInfo = () => hostelInfo;