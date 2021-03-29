CREATE TABLE IF NOT EXISTS `orders`
(
    `id` bigint(11) NOT NULL,
    `customer_id` bigint(11) NOT NULL,
    `merchant_id` bigint(11) NOT NULL,
    `amount` bigint(11) NOT NULL,
    `currency` varchar(255) NOT NULL,
    `customer_email` varchar(255) NOT NULL,
    `reference_id` bigint(11) NOT NULL,
    `status` varchar(255) NOT NULL,
    `mode` varchar(255) NOT NULL,
    `metadata` json DEFAULT NULL,
    `created_at` bigint(11) NOT NULL,
    `updated_at` bigint(11) NOT NULL,
    PRIMARY KEY (`id`)
)
DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `topups`
(
    `id` bigint(11) NOT NULL,
    `customer_id` bigint(11) NOT NULL,
    `payment_method_id` bigint(11) NOT NULL,
    `amount` bigint(11) NOT NULL,
    `currency` varchar(255) NOT NULL,
    `status` varchar(255) NOT NULL,
    `metadata` json DEFAULT NULL,
    `created_at` bigint(11) NOT NULL,
    `updated_at` bigint(11) NOT NULL,
    PRIMARY KEY (`id`)
)
DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `refunds`
(
    `id` bigint(11) NOT NULL ,
    `order_id` bigint(11) NOT NULL,
    `amount` bigint(11) NOT NULL,
    `status` varchar(255) NOT NULL,
    `refundCount` int(8) NOT NULL,
    `metadata` json DEFAULT NULL,
    `created_at` bigint(11) NOT NULL,
    `updated_at` bigint(11) NOT NULL,
    PRIMARY KEY (`id`)
)
DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `transfers`
(
    `id` bigint(11) NOT NULL,
    `order_id` bigint(11) NOT NULL,
    `type` varchar(255) NOT NULL,
    `amount` bigint(11) NOT NULL,
    `currency` varchar(255) NOT NULL,
    `status` varchar(255) NOT NULL,
    `created_at` bigint(11) NOT NULL,
    `updated_at` bigint(11) NOT NULL,
    PRIMARY KEY (`id`)
)
DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `transfer_transactions`
(
    `id` bigint(11) NOT NULL,
    `transfer_id` bigint(11) NOT NULL,
    `source_id` bigint(11) NOT NULL,
    `destination_id` bigint(11) NOT NULL,
    `wallet_pid` bigint(11) DEFAULT NULL,
    `gateway_request_id` bigint(11) DEFAULT NULL,
    `type` varchar(255) NOT NULL,
    `amount` bigint(11) NOT NULL,
    `currency` varchar(255) NOT NULL,
    `status` varchar(255) NOT NULL,
    `errorCode` varchar(255) DEFAULT NULL,
    `errorMsg` varchar(255) DEFAULT NULL,
    `metadata` json DEFAULT NULL,
    `created_at` bigint(11) NOT NULL,
    `updated_at` bigint(11) NOT NULL,
    PRIMARY KEY (`id`)
)
DEFAULT CHARSET = utf8mb4;
