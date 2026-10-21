import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('poner, ver y quitar la fecha límite', async ({ page, tasks }) => {
    await tasks([{ title: 'Llamar al proveedor' }]);
    await page.reload();

    // 1. Abre la edición y pone una fecha límite
    await page.getByRole('link', { name: 'Editar tarea: Llamar al proveedor' }).click();
    await expect(page.getByRole('heading', { name: 'Editar tarea' })).toBeVisible();
    await page.getByLabel('Fecha límite').fill('2026-10-25');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();

    // 2. La lista la muestra, también tras recargar
    await expect(page.getByRole('status')).toHaveText('Tarea actualizada.');
    await expect(page.locator('.task-due')).toHaveText('Fecha límite: 25 oct 2026');
    await page.reload();
    await expect(page.locator('.task-due time')).toHaveText('25 oct 2026');

    // 3. Vuelve a editar, vacía el campo y la fecha desaparece
    await page.getByRole('link', { name: 'Editar tarea: Llamar al proveedor' }).click();
    await expect(page.getByLabel('Fecha límite')).toHaveValue('2026-10-25');
    await page.getByLabel('Fecha límite').fill('');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();

    await expect(page.getByRole('status')).toHaveText('Tarea actualizada.');
    await expect(page.locator('.task-due')).toHaveCount(0);
  });
});
