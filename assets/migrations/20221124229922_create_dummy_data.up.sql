INSERT INTO "hubs" ("id", "name", "location", "created_at", "updated_at") OVERRIDING SYSTEM VALUE
VALUES (1, 'un assigned', 'un assigned', '0001-01-01 00:00:00', '0001-01-01 00:00:00'),
       (2, 'Shopee non tech', 'Ho Chi Minh', '2022-12-17 14:12:27', '0001-01-01 00:00:00'),
       (3, 'Lazada tech hub', 'Ho Chi Minh', '2022-12-17 14:12:42', '0001-01-01 00:00:00'),
       (4, 'Facebook VN', 'Ha Noi', '2022-12-17 14:12:55', '0001-01-01 00:00:00'),
       (5, 'Tech hub', 'Ha Noi', '0001-01-01 00:00:00', '0001-01-01 00:00:00');
SELECT setval('hubs_id_seq', (SELECT MAX(id) FROM hubs)+1);
INSERT INTO "teams" ("id", "name", "type", "hub_id", "created_at", "updated_at")  OVERRIDING SYSTEM VALUE
VALUES (1, 'un assigned', 'un assigned', 1, '0001-01-01 00:00:00', '0001-01-01 00:00:00'),
       (2, 'Manual Quality Assurance', 'qa', 2, '2022-12-17 14:16:14', '0001-01-01 00:00:00'),
       (3, 'Devops', 'sre', 2, '2022-12-17 14:16:44', '0001-01-01 00:00:00'),
       (4, 'ABC Swad', 'backend', 3, '2022-12-17 14:17:05', '0001-01-01 00:00:00'),
       (5, 'Risk Management', 'nontech', 3, '2022-12-17 14:18:14', '0001-01-01 00:00:00'),
       (6, 'DEF Squad', 'nontech', 1, '2022-12-18 10:46:04', '0001-01-01 00:00:00'),
       (7, 'Platform', 'backend', 2, '2022-12-17 14:14:55', '0001-01-01 00:00:00');
SELECT setval('teams_id_seq', (SELECT MAX(id) FROM teams)+1);