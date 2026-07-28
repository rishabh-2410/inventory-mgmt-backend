SELECT
    COALESCE(SUM(quantity) FILTER (
        WHERE movement_type = 'SALE'
          AND created_at >= NOW() - INTERVAL '30 days'
    ), 0) AS units_sold,


    COALESCE(SUM(quantity) FILTER(
        WHERE movement_type = 'RECEIVE'
        AND created_at >= NOW() - INTERVAL '30 days'
    ), 0) AS units_received,

    COALESCE(SUM(quantity) FILTER (
        WHERE movement_type = 'DAMAGE'
        AND created_at >= NOW() - INTERVAL '30 days'
    ), 0) AS units_damaged,

    COALESCE(SUM(quantity) FILTER (
        WHERE movement_type = 'RETURN'
        AND created_at >= NOW() - INTERVAL '30 days'
    ), 0) AS units_returned,

FROM stock_movements
WHERE business_id = $1