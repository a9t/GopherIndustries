package main

// GlobalProductFactory a global product factory that contains all the Products
var GlobalProductFactory = newProductFactory()

// GlobalRecipeFactory a global recipe factory that contains all the Recipes
var GlobalRecipeFactory = newRecipeFactory()

const (
	// ProductResourceCopper copper
	ProductResourceCopper int = iota
	// ProductResourceIron iron
	ProductResourceIron
	// ProductResourceStone stone
	ProductResourceStone

	// ProductProcessedCopperWire copper wire
	ProductProcessedCopperWire
	// ProductProcessedCircuitBoard circuit board
	ProductProcessedCircuitBoard

	// ProductProcessedPlate iron plate
	ProductProcessedPlate
	// ProductProcessedGear gear
	ProductProcessedGear

	// ProductStructureExtractor extractor
	ProductStructureExtractor
	// ProductStructureChest chest
	ProductStructureChest
	// ProductStructureBelt belt
	ProductStructureBelt
	// ProductStructureSplitter splitter
	ProductStructureSplitter
	// ProductStructureFactory factory
	ProductStructureFactory
	// ProductStructureUnderground underground
	ProductStructureUnderground
)

// Product generated by one of the machines in the world
type Product struct {
	name           string
	representation rune
	structure      Structure
}

// ProductFactory factory for generating all the possible Products
type ProductFactory struct {
	products        map[int]*Product
	cannonicalOrder []*Product
}

// GetProduct returns the Product identified by the product id
func (pf *ProductFactory) GetProduct(id int) *Product {
	return pf.products[id]
}

func (pf *ProductFactory) addProduct(id int, p *Product) {
	pf.products[id] = p
	pf.cannonicalOrder = append(pf.cannonicalOrder, p)
}

func newProductFactory() *ProductFactory {
	pf := new(ProductFactory)
	pf.products = make(map[int]*Product)
	pf.cannonicalOrder = make([]*Product, 0)

	pf.addProduct(ProductResourceCopper, &Product{"copper", 'c', nil})
	pf.addProduct(ProductResourceIron, &Product{"iron", 'i', nil})
	pf.addProduct(ProductResourceStone, &Product{"stone", 's', nil})

	pf.addProduct(ProductProcessedCopperWire, &Product{"wire", 'w', nil})
	pf.addProduct(ProductProcessedCircuitBoard, &Product{"circuit", 'C', nil})

	pf.addProduct(ProductProcessedPlate, &Product{"plate", 'p', nil})
	pf.addProduct(ProductProcessedGear, &Product{"gear", 'g', nil})

	pf.addProduct(ProductStructureExtractor, &Product{"extractor", 'e', NewExtractor()})
	pf.addProduct(ProductStructureChest, &Product{"chest", 'S', NewChest()})
	pf.addProduct(ProductStructureBelt, &Product{"belt", 'b', NewBelt()})
	pf.addProduct(ProductStructureSplitter, &Product{"splitter", 's', NewSplitter()})
	pf.addProduct(ProductStructureFactory, &Product{"factory", 'f', NewFactory()})
	pf.addProduct(ProductStructureUnderground, &Product{"underground", 'u', NewUnderground()})

	return pf
}

// Recipe indicates the production process required for creating a new Product
type Recipe struct {
	input           map[*Product]int
	inputOrder      []*Product
	output          *Product
	productionTicks int
}

func newRecipe(p *Product, ticks int) *Recipe {
	recipe := new(Recipe)
	recipe.input = make(map[*Product]int)
	recipe.inputOrder = make([]*Product, 0)
	recipe.output = p
	recipe.productionTicks = ticks

	return recipe
}

func (r *Recipe) addInput(p *Product, c int) {
	r.input[p] = c
	r.inputOrder = append(r.inputOrder, p)
}

// RecipeFactory stores possible Recipes
type RecipeFactory struct {
	Assembly []*Recipe
}

func newRecipeFactory() *RecipeFactory {
	rp := new(RecipeFactory)
	rp.Assembly = make([]*Recipe, 0)

	pCopper := GlobalProductFactory.GetProduct(ProductResourceCopper)
	pWire := GlobalProductFactory.GetProduct(ProductProcessedCopperWire)
	recipe := newRecipe(pWire, 100)
	recipe.addInput(pCopper, 4)
	rp.Assembly = append(rp.Assembly, recipe)

	pBoard := GlobalProductFactory.GetProduct(ProductProcessedCircuitBoard)
	recipe = newRecipe(pBoard, 200)
	recipe.addInput(pWire, 6)
	rp.Assembly = append(rp.Assembly, recipe)

	pIron := GlobalProductFactory.GetProduct(ProductResourceIron)
	pGear := GlobalProductFactory.GetProduct(ProductProcessedGear)
	recipe = newRecipe(pGear, 120)
	recipe.addInput(pIron, 4)
	rp.Assembly = append(rp.Assembly, recipe)

	pPlate := GlobalProductFactory.GetProduct(ProductProcessedPlate)
	recipe = newRecipe(pPlate, 120)
	recipe.addInput(pIron, 10)
	rp.Assembly = append(rp.Assembly, recipe)

	pStone := GlobalProductFactory.GetProduct(ProductResourceStone)
	pExtractor := GlobalProductFactory.GetProduct(ProductStructureExtractor)
	recipe = newRecipe(pExtractor, 200)
	recipe.addInput(pStone, 10)
	recipe.addInput(pPlate, 10)
	rp.Assembly = append(rp.Assembly, recipe)

	pFactory := GlobalProductFactory.GetProduct(ProductStructureFactory)
	recipe = newRecipe(pFactory, 200)
	recipe.addInput(pStone, 10)
	recipe.addInput(pPlate, 10)
	rp.Assembly = append(rp.Assembly, recipe)

	pChest := GlobalProductFactory.GetProduct(ProductStructureChest)
	recipe = newRecipe(pChest, 100)
	recipe.addInput(pPlate, 4)
	rp.Assembly = append(rp.Assembly, recipe)

	pBelt := GlobalProductFactory.GetProduct(ProductStructureBelt)
	recipe = newRecipe(pBelt, 40)
	recipe.addInput(pPlate, 1)
	recipe.addInput(pGear, 1)
	rp.Assembly = append(rp.Assembly, recipe)

	pSplitter := GlobalProductFactory.GetProduct(ProductStructureSplitter)
	recipe = newRecipe(pSplitter, 80)
	recipe.addInput(pPlate, 3)
	recipe.addInput(pGear, 3)
	rp.Assembly = append(rp.Assembly, recipe)

	pSUnderground := GlobalProductFactory.GetProduct(ProductStructureUnderground)
	recipe = newRecipe(pSUnderground, 80)
	recipe.addInput(pPlate, 5)
	recipe.addInput(pGear, 5)
	rp.Assembly = append(rp.Assembly, recipe)

	return rp
}
