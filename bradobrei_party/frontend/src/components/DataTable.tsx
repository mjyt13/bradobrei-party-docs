import type { ReactNode } from 'react'

export interface TableColumn<T> {
  key: string
  header: string
  render: (row: T) => ReactNode
}

interface DataTableProps<T> {
  caption: string
  columns: Array<TableColumn<T>>
  rows: T[]
  emptyText?: string
}

export function DataTable<T>({
  caption,
  columns,
  rows,
  emptyText = 'Данные пока отсутствуют.',
}: DataTableProps<T>) {
  return (
    <div className="table-card">
      <div className="table-scroll">
        <table className="report-table">
          <caption>{caption}</caption>
          <thead>
            <tr>
              {columns.map((column) => (
                <th key={column.key}>{column.header}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.length > 0 ? (
              rows.map((row, index) => (
                <tr key={index}>
                  {columns.map((column) => (
                    <td key={column.key}>{column.render(row)}</td>
                  ))}
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={columns.length} className="empty-state">
                  {emptyText}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
