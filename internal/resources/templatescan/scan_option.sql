SELECT base.base AS dir, file.file AS path, option.option AS option FROM file
	LEFT JOIN base ON base.id = file.base_id
	LEFT JOIN file_option ON file.id = file_option.file_id
	LEFT JOIN option ON file_option.option_id = option.id