# Есть два словаря old и new.

# Ключи — только строки. Значения — либо словари, либо строки, либо числа.

# Надо сравнить их flatten‑представления (плоский словарь, где вложенные ключи соединены точкой, например {'a': {'x': 1}} -> {'a.x': 1}) и вернуть diff в текстовом виде, как в примере:

# python
# diff = dict_diff({'a': {'x': 1}, 'b': 2},
#                  {'b': 3, 'c': 4})
# print(diff)
"""
- a.x 1
- b 2
+ b 3
+ c 4
"""


def flatten_dict(d: dict, prefix: str = "") -> dict:
    result = {}

    for key, value in d.items():
        # формируем полный ключ
        full_key = key if not prefix else f"{prefix}.{key}"

        if isinstance(value, dict):
            # рекурсивно разворачиваем вложенный словарь
            nested_flat = flatten_dict(value, full_key)
            result.update(nested_flat)
        else:
            # листовое значение
            result[full_key] = value

    return result


def dict_diff(old: dict, new: dict) -> str:
    # 1. Разворачиваем словари
    flat_old = flatten_dict(old)
    flat_new = flatten_dict(new)

    # 2. Собираем все ключи
    all_keys = sorted(set(flat_old.keys()) | set(flat_new.keys()))

    lines = []

    # 3. Для каждого ключа смотрим, что произошло
    for key in all_keys:
        in_old = key in flat_old
        in_new = key in flat_new

        if in_old and not in_new:
            # ключ исчез
            lines.append(f"- {key} {flat_old[key]}")
        elif not in_old and in_new:
            # ключ появился
            lines.append(f"+ {key} {flat_new[key]}")
        else:
            # есть в обоих
            old_val = flat_old[key]
            new_val = flat_new[key]
            if old_val != new_val:
                # изменилось значение: сначала - старое, потом + новое
                lines.append(f"- {key} {old_val}")
                lines.append(f"+ {key} {new_val}")

    # 4. Собираем итоговую строку
    return "\n".join(lines)


if __name__ == "__main__":
    diff = dict_diff({"a": {"x": 1}, "b": 2}, {"b": 3, "c": 4})
    print(diff)
