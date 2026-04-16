-- Default categories (user_id = NULL = shared/system, không thuộc về user nào)
INSERT INTO categories (user_id, name, color, icon, is_default, sort_order) VALUES
    (NULL, 'Làm việc',    '#3B82F6', 'briefcase',  TRUE, 1),
    (NULL, 'Di chuyển',   '#F59E0B', 'car',         TRUE, 2),
    (NULL, 'Ăn uống',     '#10B981', 'utensils',    TRUE, 3),
    (NULL, 'Giải trí',    '#8B5CF6', 'gamepad',     TRUE, 4),
    (NULL, 'Mạng xã hội', '#EF4444', 'smartphone',  TRUE, 5),
    (NULL, 'Ngủ nghỉ',    '#6B7280', 'moon',        TRUE, 6);
