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
    `updated_at` bigint(11) NOT NULL
)
DEFAULT CHARSET = utf8mb4;
