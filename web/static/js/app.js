const ORDER_STATUSES = {
    1: { text: "In Store", color: "blue" },
    2: { text: "In Transit", color: "orange" }, 
    3: { text: "Delivered", color: "green" },
    4: { text: "Processing", color: "purple" }
};

async function loadOrder() {
    const orderId = document.getElementById('orderIdInput').value.trim();
    const container = document.getElementById('order-container');
    
    if (!orderId) {
        container.innerHTML = '<div class="error-message">Please enter an Order ID</div>';
        return;
    }

    container.innerHTML = `
        <div class="loading">
            <div class="loading-spinner"></div>
            <div>Loading order ${orderId}...</div>
        </div>
    `;

    try {
        const response = await fetch(`/orders/${orderId}`);
        if (!response.ok) {
            throw new Error(`Order not found (Status: ${response.status})`);
        }
        
        const order = await response.json();
        renderOrder(order);
    } catch (error) {
        container.innerHTML = `
            <div class="error-message">
                <strong>Error:</strong> ${error.message}
            </div>
        `;
    }
}

function renderOrder(order) {
    const container = document.getElementById('order-container');
    
    const orderStatus = determineOrderStatus(order);
    const statusConfig = ORDER_STATUSES[orderStatus];
    
    container.innerHTML = `
        <div class="order-card">
            <div class="card-header">
                <div class="order-header-with-status">
                    <div class="order-basic-info">
                        <div class="order-id">Order #${order.order_uid}</div>
                        <div class="order-meta">
                            <span>📦 ${order.track_number}</span>
                            <span>🏢 ${order.entry}</span>
                            <span>🕐 ${new Date(order.date_created).toLocaleString()}</span>
                        </div>
                    </div>
                    <div class="order-status-section">
                        <div class="order-total-amount">$${(order.payment.amount / 100).toFixed(2)}</div>
                        <div class="order-status-badge status-${statusConfig.color}">
                            ${statusConfig.text}
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="card-body">
                <div class="section">
                    <div class="section-title">📋 Order Information</div>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="info-label">Customer ID</span>
                            <span class="info-value">${order.customer_id}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Delivery Service</span>
                            <span class="info-value">${order.delivery_service}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Locale</span>
                            <span class="info-value">${order.locale.toUpperCase()}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Internal Signature</span>
                            <span class="info-value">${order.internal_signature || 'N/A'}</span>
                        </div>
                    </div>
                </div>

                <div class="section">
                    <div class="section-title">🚚 Delivery Details</div>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="info-label">Recipient</span>
                            <span class="info-value">${order.delivery.name}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Phone</span>
                            <span class="info-value">${order.delivery.phone}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Email</span>
                            <span class="info-value">${order.delivery.email}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Address</span>
                            <span class="info-value">${order.delivery.city}, ${order.delivery.address}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">ZIP Code</span>
                            <span class="info-value">${order.delivery.zip}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Region</span>
                            <span class="info-value">${order.delivery.region}</span>
                        </div>
                    </div>
                </div>

                <div class="section">
                    <div class="section-title">💳 Payment Information</div>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="info-label">Transaction ID</span>
                            <span class="info-value">${order.payment.transaction}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Total Amount</span>
                            <span class="info-value">$${(order.payment.amount / 100).toFixed(2)} ${order.payment.currency}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Provider</span>
                            <span class="info-value">${order.payment.provider}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Bank</span>
                            <span class="info-value">${order.payment.bank}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Delivery Cost</span>
                            <span class="info-value">$${(order.payment.delivery_cost / 100).toFixed(2)}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Goods Total</span>
                            <span class="info-value">$${(order.payment.goods_total / 100).toFixed(2)}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Custom Fee</span>
                            <span class="info-value">$${(order.payment.custom_fee / 100).toFixed(2)}</span>
                            <small style="color: #6b7280; font-size: 0.8em;">(Additional customs charges)</small>
                        </div>
                    </div>
                </div>

                <div class="section">
                    <div class="section-title">🛍️ Order Items (${order.items.length})</div>
                    <div class="product-items-grid">
                        ${order.items.map(item => {
                            const originalPrice = (item.price / 100) * item.quantity;
                            const discountAmount = originalPrice * (item.sale / 100);
                            const finalPrice = originalPrice - discountAmount;
                            
                            return `
                                <div class="product-card">
                                    <div class="product-header">
                                        <div class="product-name">${item.name}</div>
                                        <div class="product-quantity">×${item.quantity}</div>
                                    </div>
                                    
                                    <div class="product-content">
                                        <div class="product-price-section">
                                            <div class="price-grid">
                                                <div class="price-row">
                                                    <span class="price-label">Unit Price</span>
                                                    <span class="price-value">$${(item.price / 100).toFixed(2)}</span>
                                                </div>
                                                <div class="price-row">
                                                    <span class="price-label">Quantity</span>
                                                    <span class="price-value">${item.quantity}</span>
                                                </div>
                                                <div class="price-row">
                                                    <span class="price-label">Subtotal</span>
                                                    <span class="price-value">$${originalPrice.toFixed(2)}</span>
                                                </div>
                                                <div class="price-row">
                                                    <span class="price-label">Discount (${item.sale}%)</span>
                                                    <span class="price-value price-discount">-$${discountAmount.toFixed(2)}</span>
                                                </div>
                                            </div>
                                            <div class="price-row price-total">
                                                <span class="price-label">Item Total</span>
                                                <span class="price-value">$${finalPrice.toFixed(2)}</span>
                                            </div>
                                        </div>
                                        
                                        <div class="product-details-section">
                                            <div class="details-title">Product Details</div>
                                            <div class="details-grid">
                                                <div class="detail-card">
                                                    <div class="detail-label">Size</div>
                                                    <div class="detail-value">${item.size}</div>
                                                </div>
                                                <div class="detail-card">
                                                    <div class="detail-label">Brand</div>
                                                    <div class="detail-value">${item.brand}</div>
                                                </div>
                                                <div class="detail-card">
                                                    <div class="detail-label">Sale</div>
                                                    <div class="detail-value detail-sale">${item.sale}% OFF</div>
                                                </div>
                                                <div class="detail-card">
                                                    <div class="detail-label">Status</div>
                                                    <div class="detail-value detail-status status-${getItemStatusClass(item.status)}">${item.status}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            `;
                        }).join('')}
                    </div>
                </div>
            </div>
        </div>
    `;
}

function determineOrderStatus(order) {
    if (order.status) {
        return order.status;
    }
    
    if (order.items.every(item => item.status === "delivered")) return 3;
    if (order.items.some(item => item.status === "shipped")) return 2;
    return 1;
}

function getItemStatusClass(status) {
    const statusMap = {
        'pending': 'pending',
        'processing': 'processing', 
        'shipped': 'shipped',
        'delivered': 'delivered'
    };
    return statusMap[status] || 'pending';
}

document.addEventListener('DOMContentLoaded', function() {
    const orderInput = document.getElementById('orderIdInput');
    if (orderInput) {
        orderInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                loadOrder();
            }
        });
    }
});

function loadSpecificOrder(orderId) {
    document.getElementById('orderIdInput').value = orderId;
    loadOrder();
}