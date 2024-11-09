<template>    
    <nav>
    <h1>Stickerpack editor</h1>
        <ul>
            <li>
                <router-link 
                    :to="`/my-stickerpacks`"
                    class="link"
                    active-class="active"
                ><span>My stickerpacks</span></router-link>
            </li>
        </ul>
    </nav>
</template>

<style scoped>
h1 {
    font-size: 1.5em;
    text-align: center;
    padding: 20px 40px;
}

nav {
    display: flex;
    background: var(--secondary);
    padding: 1px;
}

ul {
    display: flex;
    flex: 1;
    justify-content: center;
    margin-left: -336px;
}

li {
    display: flex;
}

.link {
    padding: 20px;
    text-decoration: none;
    color: var(--text);
    font-size: 1.2em;
    border-left: 3px solid var(--secondary);
    border-right: 3px solid var(--secondary);
    display: flex;
    justify-content: center;
    align-items: center;
}

.link:hover {
    background: var(--accent);
    border-color: var(--accent)
}

.active {
    border-color: var(--text)
}

.active:hover {
    border-color: var(--text);
}
</style>