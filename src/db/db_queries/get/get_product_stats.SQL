SELECT
    COALESCE(
            SUM(
                CASE
                    WHEN movement_type='SALE'
                    THEN movement_quantity 
                END
            ), 
            0  
    ) AS total_sold,
    COALESCE(
            SUM(
                CASE
                    WHEN movement_type='RETURN'
                    THEN movement_quantity 
                END
            ), 
            0
    ) AS total_return,
    COALESCE(
            SUM(
                CASE
                    WHEN movement_type='DAMAGE'
                    THEN movement_quantity 
                END
            ), 
            0
    ) AS total_damage,
    COALESCE(
            SUM(
                CASE
                    WHEN movement_type='RECEIVE'
                    THEN movement_quantity 
                END
            ), 
            0
    ) AS total_receive
FROM stock_movement
WHERE business_id='9625c084-99b2-4119-920a-eabd9bca6e15' AND product_id='7388bc7e-13b9-4165-840b-b0189536c00b'
AND created_at >= NOW() - INTERVAL '3 months'